/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package main

func build_gtk_func() {

	special_func["draw_init"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
		go run_draw_app()

		special_func["draw_clear"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			draw_clear()
			return NewNil(), nil
		}
		special_func["draw_line"] = func(env *SimpleEnv, v []Expression) (Expression, error) {
			var point [4]int
			for i, sexp := range v {
				e, err := eval(sexp, env)
				if err != nil {
					return nil, err
				}
				if p, ok := e.(*Integer); ok {
					point[i] = p.Value
				} else if p, ok := e.(*Float); ok {
					point[i] = int(p.Value)
				} else {
					NewRuntimeError("E1003", reflect.TypeOf(e).String())
				}
			}
			draw_line_reentrant(point[0], point[1], point[2], point[3])
			return NewNil(), nil
		}
		return NewNil(), nil
	}
}
// Main
func main() {

	build_func()
	build_gtk_func()

	cui_ch := make(chan bool)
	go func() {
		do_interactive()
		cui_ch <- true
	}()
	<-cui_ch
}
