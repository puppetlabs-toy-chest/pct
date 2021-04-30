package ruby

import (
	"io/fs"
	"testing"
)

// // Mock response of os.Stat(commandPath)
// type testFile struct{}

// func(testFile) Stat(string) (os.FileInfo, error) { return nil, nil}
// func(testFile) Read([]byte) (int, error){return 0, nil}
// func(testFile) Close() error {return nil}

// Mock response of utils.ValidModuleRoot()
// type tUtils struct {}

// func (t *tUtils) Contains([]string, string) bool {return true}
// func (t *tUtils) Find([]string, string) ([]string) {return []string{}}
// func (t *tUtils) GetListOfFlags(cobra.Command, []string) []string {return []string{}}
// func (t *tUtils) FlagsToIgnore() []string { return []string{}}
// func (t *tUtils) ExecutePDKCommand(cmd *cobra.Command, args []string) error {return nil}
// func (t *tUtils) ValidModuleRoot() (moduleRoot string, err error) {return "/foo/bar/", nil}


func TestRubyRubocopValidator_SetCommand(t *testing.T) {
	type args struct {
		command string
		options []string
	}

	var tests[] struct {
		name    string
		args    args
		wantErr bool
	}

	// origValidModuleRoot := fValidModuleRoot
	// defer func() { fValidModuleRoot = origValidModuleRoot }()
	fValidModuleRoot = func() (string, error) {return "/foo/bar", nil}

	// origfOsStat := fOsStat
	// defer func() { fOsStat = origfOsStat }()
	fOsStat = func(string) (fs.FileInfo, error) {return nil, nil}

	tests = append(tests, struct {
		name 		string
		args		args
		wantErr bool
	}{
		name: 		"Sets command file path and args",
		args: 		args{
			command: "rubocop",
			options: []string {
				"foo",
				"bar",
			},
		},
		wantErr: 	false,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RubyRubocopValidator{}
			if err := r.SetCommand(tt.args.command, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("RubyRubocopValidator.SetCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
