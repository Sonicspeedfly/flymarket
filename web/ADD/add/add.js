async function scorGame() {
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1);
    result = score[0]
    list = document.querySelector(".addproduct")
    list.innerHTML += `
    <form action="http://127.0.0.1:9999/api/products/save?id=0&accountid=${result}" method="POST" enctype="multipart/form-data">
    <input type="text" placeholder="Название продукта" name="name">
    <input type="text" placeholder="Информация" name="information">
    <select name="category">
    <option value="smartphone">Смартфон</option>
    <option value="computer">Компьютер</option>
    <option value="TV">Телевизор</option>
    </select>
    <input type="number" placeholder="Количество" name="count">
    <input type="number" placeholder="Цена" name="price">
    <h3>Фото продукта:</h3>
    <input type="file" name="image">
    <button>Отправить</button>
    </form>
    `
}