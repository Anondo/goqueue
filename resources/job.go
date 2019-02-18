package resources

type Job struct {
	ID      int
	JobName string        `json:"job_name"`
	Args    []interface{} `json:"args"`
}

// func (j *Job) Execute() {
// 	j.F.Call(j.Args)
// }

//
// func (j *Job) ProcessJob(f interface{}, args ...interface{}) error {
// 	j.F = reflect.ValueOf(f)
// 	j.FType = reflect.TypeOf(f)
//
// 	if j.FType.NumIn() != len(args) {
// 		return errors.New("Invalid Number Of Arguments For Job")
// 	}
//
// 	for i, arg := range args {
// 		if reflect.TypeOf(arg) != j.FType.In(i) {
// 			return fmt.Errorf("Argmunet Type Mismatch: Expected %v, Got %v", j.FType.In(i), reflect.ValueOf(arg))
// 		}
// 		j.Args = append(j.Args, reflect.ValueOf(arg))
// 	}
//
// 	return nil
// }
