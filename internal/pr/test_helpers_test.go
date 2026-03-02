package pr

import "fmt"

type fakeRunner struct {
	responses map[string]string
	errors    map[string]error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{responses: map[string]string{}, errors: map[string]error{}}
}

func (f *fakeRunner) when(args []string, response string) {
	f.responses[key(args)] = response
}

func (f *fakeRunner) whenErr(args []string, err error) {
	f.errors[key(args)] = err
}

func (f *fakeRunner) Run(args ...string) (string, error) {
	k := key(args)
	if err, ok := f.errors[k]; ok {
		return "", err
	}
	if out, ok := f.responses[k]; ok {
		return out, nil
	}
	return "", fmt.Errorf("unexpected gh args: %v", args)
}

func key(args []string) string {
	return fmt.Sprintf("%q", args)
}
