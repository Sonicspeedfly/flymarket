async function scorGame() {
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1);
    result = score[0]
    list = document.querySelector(".addproduct")
    list.innerHTML += `
    <form action="http://127.0.0.1:9999/api/products/save?id=0&accountid=${result}" method="POST" enctype="multipart/form-data">
    <input type="text" name="name">
    <input type="text" name="information">
    <input type="number" name="count">
    <input type="file" name="image">
    <button>Отправить</button>
    </form>
    `
}