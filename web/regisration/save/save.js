async function register() {
    phone = document.getElementById("form3Example3cg").value;
    email = document.getElementById("form3Example2cg").value;
    username = document.getElementById("form3Example1cg").value;
    password = document.getElementById("form3Example4cg").value;
    url = 'http://127.0.0.1:9999/api/accounts?phone=' + phone +'&email=' + email + '&username=' + username + '&password=' + password;
    response = await fetch(url);
    content = await response.json()
    //console.log(content)
    window.location.replace('../login/index.html');
    // list = document.querySelector(".success")
    // if (content.active === true) {
        // list.innerHTML += `<p classs="complete"> Вы успешно зарегестрированны!</p>`
    // }
}
