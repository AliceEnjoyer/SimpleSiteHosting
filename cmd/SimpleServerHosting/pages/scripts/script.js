const MsgQAcceptCookies = "This website uses cookies in order to offer you the most relevant information. Please accept cookies for optimal performance."

if(localStorage.getItem("isCookies").toString() == "false" 
    && confirm(MsgQAcceptCookies)){
    alert("Now we can offer you the most relevant information!")
    localStorage.setItem("isCookies", "true")
}
