package test

import (
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"github.com/facebookgo/inject"
	"testing"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"reflect"
	"os"
)

func TestPersistenceManager_RecoverLeases(t *testing.T) {


	fixtures := []struct {
		fileUtils        FakeFileUtils
		expectedError error
	}{
		{
			fileUtils: FakeFileUtils{
				readEnv: DefaultReadEnvMetrics(),
				readDir: ReadDirTestMetrics{
					content: []os.FileInfo{
						os.FileInfo(FakeFileInfo{"test"}),
					},
					MethodCallMetrics: MethodCallMetrics{
						expectedCalls: 1,
						method: "readdir",
					},
				},

				readFile: ReadFileTestMetrics{
					content: "{}",
					MethodCallMetrics: MethodCallMetrics{
						expectedCalls: 1,
						method: "readfile",
					},
				},

			},
			expectedError: nil,
		},
	}

	for i, fixture := range fixtures {
		cfg := &config.Config{}

		if err := inject.Populate(cfg, &fixture.fileUtils); err != nil {
			t.Error(err.Error())
		}

		cfg.ReadConfig()

		persistenceCfg := &persistence.PersistenceManager{}
		if err := inject.Populate(persistenceCfg,&fixture.fileUtils); err != nil {
			t.Error(err.Error())
		}
		_, persistenceManager := persistence.NewPersistenceManager(cfg,persistenceCfg)

		err := persistenceManager.RecoverLeases()
		fixture.fileUtils.readEnv.Report(t, i)
		fixture.fileUtils.readFile.Report(t, i)
		fixture.fileUtils.readDir.Report(t, i)

		if !reflect.DeepEqual(err, fixture.expectedError) {
			t.Errorf("%d - Expected %v, Received %v\n", i, fixture.expectedError, err)
		}
	}
}
