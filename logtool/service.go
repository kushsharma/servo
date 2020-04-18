package logtool

// Service that implements logtool.Service interface
type Service struct {
}

// Fetch extracts the contents of file
func (svc *Service) Fetch(path, filename string) (string, error) {
	content := ""

	return content, nil
}

//List return file names in the directory
func (svc *Service) List(path string) ([]string, error) {
	files := []string{}

	return files, nil
}

func (svc *Service) Delete(path string) error {

	return nil
}

func NewService() *Service {
	return &Service{}
}
