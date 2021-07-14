async function register() {
    phone = document.getElementById("phone").value;
    email = document.getElementById("email").value;
    username = document.getElementById("username").value;
    password = document.getElementById("password").value;
 
    url = 'http://127.0.0.1:9999/api/accounts?phone=' + phone +'&email=' + email + '&username=' + username + '&password=' + password;
    response = await fetch(url);
    content = await response.json()
    console.log(content)
    list = document.querySelector(".success")
    if (content.active === true) {
        list.innerHTML += `<p classs="complete"> Вы успешно зарегестрированны!</p>`
    }
}
