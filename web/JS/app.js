async function getResponse() {
    response = await fetch('http://127.0.0.1:9999/api/products/All');
    content = await response.json()
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1, 2);
    var reg = score[2];

    list = document.querySelector('.products')

    let key;

    if (reg === "true") {
        for (key in content) {
        
            list.innerHTML += `
                <dl class="product">
                    <h4>${content[key].name}</h4>
                    <img class="imgproduct" src="./banners/${content[key].id}/${content[key].file}">
                    <p class="textproduct">${content[key].information}</p>
                    <p>количество: ${content[key].count}</p>
                    <script src="./JS/app.js"></script>
                    <input type="number" id="count${content[key].id}">
                    <button onclick="buy(${content[key].id})">Купить</button> 
                    <div class="con${content[key].id}">
                    </div> 
                </dl>
                    `
            
        }
    };
    if (reg === undefined) {
        for (key in content) {
        
        list.innerHTML += `
            <dl class="product">
                <h4>${content[key].name}</h4>
                <img class="imgproduct" src="./banners/${content[key].id}/${content[key].file}">
                <p class="textproduct">${content[key].information}</p>
                <p>количество: ${content[key].count}</p>
            </dl>
                `
    }
}
}

async function scoreGame() {
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1, 2);
    var result = score[1];
    var reg = score[2];
    list = document.querySelector(".user")
    if (reg === "true") {
        responses = await fetch('http://127.0.0.1:9999/api/accounts/autification?token=' + result);
            content = await responses.json()
            list.innerHTML = `<p class="login">${content.username}</p>
            <button class="add"><a href="./ADD/index.html?&${content.id}">Добавить</a></button>`
            getResponse()
    };
    if (reg === undefined) {
        list.innerHTML += `
        <a class="reg" href="./regisration/index.html">Регистрация</a>
        <a class="login" href="./login/index.html">Войти</a>
        `
        getResponse()
    }
    return result, reg
}

async function buy(id) {
    var score = decodeURIComponent(location.search.substr(1)).split('&');
    score.splice(0, 1, 2);
    var result = score[1];
    var reg = score[2];
    count = document.getElementById("count" + id).value;
    responses = await fetch('http://127.0.0.1:9999/api/products/Buy?id=' + id + '&count=' + count);
    content = await responses.json();
    window.open('./index.html?&' + result + '&' + reg)
    lister = document.querySelector(".con"+id);
    lister.innerHTML += `
    <p>плокупка произведена успешно!</p> 
    `
}