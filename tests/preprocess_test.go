package tests

import (
	"testing"

	"github.com/HexmosTech/lama2/preprocess"
	"github.com/rs/zerolog/log"
)

func TestPreprocessBasic(t *testing.T) {
	op := preprocess.PreprocessLamaFile("../elfparser/ElfTestSuite/env1/sample.l2")
	log.Debug().Str("Preprocessed string", op)
}
