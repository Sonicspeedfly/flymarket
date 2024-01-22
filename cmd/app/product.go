package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sonicspeedfly/flymarket/pkg/product"
)

func (s *Server) handleSaveProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	err := request.ParseMultipartForm(64 << 40)
	if err != nil {
		log.Println(err)
		return
	}
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	accidParam := request.URL.Query().Get("accountid")
	nameParam := request.FormValue("name")
	imageParam := request.FormValue("image")
	categoryParam := request.FormValue("category")
	informationParam := request.FormValue("information")
	countParam := request.PostFormValue("count")
	priceParam := request.FormValue("price")
	images := filepath.Ext(imageParam)

	file, fileHeader, err := request.FormFile("image")

	if err == nil {
		var name = strings.Split(fileHeader.Filename, ".")
		images = name[len(name)-1]
	}

	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Println(accidParam)
	accid, err := strconv.ParseInt(accidParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseInt(priceParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.productSvc.SaveProduct(request.Context(), id, nameParam, images, categoryParam, nameParam, informationParam, count, price, accid)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	namefile := item.NameProduct + "." + images
	idtext := strconv.Itoa(int(item.ID))
	uploadFile(file, idtext, "./web/banners/", namefile)
	item, _ = s.productSvc.SaveProduct(request.Context(), item.ID, item.NameProduct, item.Image, item.Category, namefile, item.Information, item.Count, item.Price, item.Account_ID)
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleEditProduct(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		return
	}

	idParam := request.PostFormValue("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	nameParam := request.FormValue("name")
	imageParam := request.FormValue("image")
	categoryParam := request.FormValue("category")
	informationParam := request.FormValue("information")
	countParam := request.PostFormValue("count")
	priceParam := request.PostFormValue("price")
	accidParam := request.PostFormValue("accountid")
	images := filepath.Ext(imageParam)

	file, fileHeader, err := request.FormFile("image")

	if err == nil {
		var name = strings.Split(fileHeader.Filename, ".")
		images = name[len(name)-1]
	}

	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	accid, err := strconv.ParseInt(accidParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseInt(priceParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.productSvc.SaveProduct(request.Context(), id, nameParam, images, categoryParam, nameParam, informationParam, count, price, accid)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	namefile := item.NameProduct + "." + images
	idtext := strconv.Itoa(int(item.ID))
	uploadFile(file, idtext, "../web/banners/", namefile)
	item, _ = s.productSvc.SaveProduct(request.Context(), item.ID, item.NameProduct, item.Image, item.Category, namefile, item.Information, item.Count, item.Price, item.Account_ID)
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleByIdProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.productSvc.ByID(request.Context(), id)
	if errors.Is(err, errors.New("Not Found")) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleAllProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	item, err := s.productSvc.All(request.Context())
	if errors.Is(err, errors.New("Not Found")) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveProduct(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.productSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleBuyProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	countParam := request.URL.Query().Get("count")

	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.productSvc.BuyProduct(request.Context(), id, count)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleInstallment(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	idParam := request.URL.Query().Get("id")
	monthsParam := request.URL.Query().Get("months")
	categoryParam := request.URL.Query().Get("category")
	priceParam := request.URL.Query().Get("price")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	months, err := strconv.ParseInt(monthsParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseInt(priceParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := product.Calc(request.Context(), id, months, categoryParam, price)
	if errors.Is(err, errors.New("Not Found")) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func uploadFile(file multipart.File, dir string, path string, namefile string) error {
	err := os.MkdirAll(path+dir, 0777)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path+dir+"/"+namefile, data, 0666)

	if err != nil {
		return errors.New("not saved from folder ")
	}
	return nil
}
