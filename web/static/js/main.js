// ── Boutons de vote (like / dislike) ─────────────────────────────────────

document.querySelectorAll(".react-btn").forEach(function(bouton) {
  bouton.addEventListener("click", function() {
    var id    = bouton.dataset.id;
    var type  = bouton.dataset.type;
    var value = bouton.dataset.value;

    fetch("/react", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: "target_id=" + id + "&target_type=" + type + "&value=" + value
    })
    .then(function(reponse) {
      return reponse.json();
    })
    .then(function(data) {
      // on cherche la barre de réaction qui contient ce bouton
      var barre = bouton.closest(".reaction-bar");
      if (!barre) return;

      var btnLike    = barre.querySelector(".like-btn");
      var btnDislike = barre.querySelector(".dislike-btn");
      if (!btnLike || !btnDislike) return;

      // mise à jour des compteurs
      btnLike.querySelector(".react-count").textContent    = data.likes;
      btnDislike.querySelector(".react-count").textContent = data.dislikes;

      // mise à jour des états actif / inactif
      btnLike.classList.toggle("active",    data.user_reaction === 1);
      btnDislike.classList.toggle("active", data.user_reaction === -1);

      // mise à jour aria-pressed pour accessibilité
      btnLike.setAttribute("aria-pressed",    String(data.user_reaction === 1));
      btnDislike.setAttribute("aria-pressed", String(data.user_reaction === -1));

      // petite animation
      bouton.classList.add("bumping");
      bouton.addEventListener("animationend", function() {
        bouton.classList.remove("bumping");
      }, { once: true });
    })
    .catch(function(erreur) {
      console.error("Erreur lors du vote :", erreur);
    });
  });
});

// ── Audio description (lecture du post à voix haute) ─────────────────────

var btnAudio = document.getElementById("btn-audio");

if (btnAudio && window.speechSynthesis) {
  var enLecture = false;

  btnAudio.addEventListener("click", function() {
    var synth = window.speechSynthesis;

    // si déjà en lecture → on arrête
    if (enLecture) {
      synth.cancel();
      enLecture = false;
      btnAudio.classList.remove("en-lecture");
      btnAudio.querySelector("#audio-label").textContent = "Écouter";
      btnAudio.setAttribute("aria-label", "Lire ce post à voix haute");
      return;
    }

    // on récupère le titre et le contenu
    var titre   = btnAudio.dataset.titre || "";
    var contenuId = btnAudio.dataset.contenuId;
    var contenu = "";
    if (contenuId) {
      var element = document.getElementById(contenuId);
      if (element) contenu = element.textContent || "";
    }

    var texteALire = titre + ". " + contenu;

    var utterance = new SpeechSynthesisUtterance(texteALire);
    utterance.lang = "fr-FR";
    utterance.rate = 0.95;
    utterance.pitch = 1;

    // quand la lecture se termine tout seule
    utterance.onend = function() {
      enLecture = false;
      btnAudio.classList.remove("en-lecture");
      btnAudio.querySelector("#audio-label").textContent = "Écouter";
      btnAudio.setAttribute("aria-label", "Lire ce post à voix haute");
    };

    synth.speak(utterance);
    enLecture = true;
    btnAudio.classList.add("en-lecture");
    btnAudio.querySelector("#audio-label").textContent = "Stop";
    btnAudio.setAttribute("aria-label", "Arrêter la lecture");
  });

} else if (btnAudio) {
  // le navigateur ne supporte pas la synthèse vocale
  btnAudio.disabled = true;
  btnAudio.title = "Votre navigateur ne supporte pas la lecture audio";
}
