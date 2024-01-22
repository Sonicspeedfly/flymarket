async function loginin() {
    login = document.getElementById("form1Example13").value;
    password = document.getElementById("form1Example23").value;
 
    //alert(login + ' ' + password);
    url = 'http://127.0.0.1:9999/api/accounts/login?phone=' + login +'&password=' + password;
    response = await fetch(url);
    content = await response.json();
    //console.log(content)
    document.cookie = 'token=' + content.token + '; path=/web/; domain=127.0.0.1';
    document.cookie ='reg=true; path=/web/; domain=127.0.0.1';
    //alert(getCookie('token'));
    window.location.replace('../index.html?');
}
