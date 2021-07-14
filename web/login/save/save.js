async function loginin() {
    login = document.getElementById("phone").value;
    password = document.getElementById("password").value;
 
    url = 'http://127.0.0.1:9999/api/accounts/login?phone=' + login +'&password=' + password;
    response = await fetch(url);
    content = await response.json()
    console.log(content)
    window.open('../index.html?&'+content.token+'&true');
}
