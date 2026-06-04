package forumjs

import "net/http"

func ChangerDeFilDeDiscution(w http.ResponseWriter, r *http.Request) {
	RevenirSurLaPageAccueil(w,r,0,false,true,-1)
}
