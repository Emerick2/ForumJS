function Partager(){
    const url = window.location.href;

    navigator.clipboard.writeText(url).then(() => {
        alert("L'URL de la page à été copié !")
    })
}