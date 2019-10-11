package loaders

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoader(t *testing.T) {

	testFile := "test/test-resource.txt"
	fileContent := "PRODUCTION GRID RESOURCE LOADER TEST FILE\n"

	assert := assert.New(t)

	loader := FileResourceLoader{}

	//test main reader method
	reader, err := loader.Reader(testFile)

	assert.NoError(err)
	assert.NotNil(reader)

	content, err := ioutil.ReadAll(reader)
	assert.NoError(err)

	assert.Equal(fileContent, string(content))

	//test convenience methods
	content, err = loader.Bytes(testFile)
	assert.NoError(err)
	assert.Equal(fileContent, string(content))

	text, err := loader.String(testFile)
	assert.NoError(err)
	assert.Equal(fileContent, text)

}
