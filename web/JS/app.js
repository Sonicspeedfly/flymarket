function get_cookie ( cookie_name )
{
  var results = document.cookie.match ( '(^|;) ?' + cookie_name + '=([^;]*)(;|$)' );
 
  if ( results )
    return ( unescape ( results[2] ) );
  else
    return null;
}

async function getResponse() {
    response = await fetch('http://127.0.0.1:9999/api/products/All');
    content = await response.json()
    var reg = get_cookie('reg');

    list = document.querySelector('.products');

    let key;

    if (reg === "true") {
        for (key in content) {
        
            list.innerHTML += `
            <div class="card mb-3" style="max-width: 100%;">
            <div class="row g-0" style="margin-top: 1%;margin-bottom:1%">
              <div class="col-md-2">
                <img src="./banners/${content[key].id}/${content[key].file}" style="max-height:240px; margin-left:5%;" class="img-fluid rounded-start" alt="${content[key].file}">
              </div>
              <div class="col-md-9" style="margin-left:1%">
                <div class="card-body">
                  <h5 class="card-title">${content[key].name}</h5>
                  <p class="card-text">${content[key].information}</p>
                  <p class="card-text"><small class="text-body-secondary">количество: ${content[key].count}</small></p>
                  <p class="card-text"><small class="text-body-secondary">ЦЕНА: ${content[key].price}</small></p>
                  <div class="input-group mb-3">
                    <input class="form-control me-0" min="1" max="${content[key].count}" type="number" id="count${content[key].id}" aria-describedby="button-addon2">
                    <button class="btn btn-outline-success" onclick="buy(${content[key].id})" id="button-addon2">Купить</button>
                  </div>
                  <div class="con${content[key].id}">
                </div>
              </div>
            </div>
          </div>
        `
            
        }
    } else {
        for (key in content) { 
        list.innerHTML += `
        <div class="card mb-3" style="max-width: 100%;">
          <div class="row g-0" style="margin-top: 1%;margin-bottom:1%">
            <div class="col-md-2">
              <img src="./banners/${content[key].id}/${content[key].file}" style="max-height:240px; margin-left:5%;" class="img-fluid rounded-start" alt="${content[key].file}">
            </div>
            <div class="col-md-9" style="margin-left:1%">
              <div class="card-body">
                <h5 class="card-title">${content[key].name}</h5>
                <p class="card-text">${content[key].information}</p>
                <p class="card-text"><small class="text-body-secondary"><b>ЦЕНА:</b> ${content[key].price}</small></p>
              </div>
            </div>
          </div>
        </div>
        `
    }
}
}

async function scoreGame() {
    var result = get_cookie('token');
    var reg = get_cookie('reg');
    console.log(result + ' ' + reg);
    list = document.querySelector(".user")
    if (reg === "true") {
        responses = await fetch('http://127.0.0.1:9999/api/accounts/autification?token=' + result);
            content = await responses.json()
            list.innerHTML = `<p class="login" style="color: rgb(3, 95, 201); font-family: rubik; font-size: 20pt;">${content.username}</p>`
            document.querySelector(".tov").innerHTML += `<button class="btn add hidden btn-success reg" onclick="document.location='./ADD/index.html?&${content.id}'" type="submit">Добавить</button>`;
            getResponse()
    } else {
        list.innerHTML += `
        <button onclick="document.location='./regisration/index.html'" class="btn btn-outline-success reg" type="submit">Регистрация</button>
        <button onclick="document.location='./login/index.html'" class="btn btn-outline-success login" type="submit">Войти</button>
        `
        getResponse()
    }
    return result, reg
}

async function buy(id) {
    var result = get_cookie('token');
    var reg = get_cookie('reg');
    count = document.getElementById("count" + id).value;
    response = await fetch('http://127.0.0.1:9999/api/products/product.ByID?id=' + id);
    conten = await response.json()
    if(count >= 1 && count <= conten.count) {
        responses = await fetch('http://127.0.0.1:9999/api/products/Buy?id=' + id + '&count=' + count);
        content = await responses.json();
        //window.open('./index.html?&' + result + '&' + reg)
        lister = document.querySelector(".con"+id);
        lister.innerHTML += `
        <p>покупка произведена успешно!</p> 
        `
    }
    
}

window.addEventListener('scroll', function() {
    if(pageYOffset>30){
        document.querySelector(".navig").classList.remove("hidden");
    } else {
        document.querySelector(".navig").classList.add("hidden");
    }
});

function linkClick() {
    pageYOffset-=42;
}