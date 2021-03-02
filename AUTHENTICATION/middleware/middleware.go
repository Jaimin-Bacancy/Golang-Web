//AuthenticationCustomer is..
func AuthenticationCustomer(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !ok {

			fmt.Println("not exits")
			http.Redirect(w, r, "/Login", http.StatusSeeOther)
			return
		}
		handler.ServeHTTP(w, r)
	}
}