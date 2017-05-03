package writer_test

import (
	. "github.com/monkeyherder/salus/checks/writer"

	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cznic/fileutil"
	"github.com/monkeyherder/salus/checks"
	"github.com/monkeyherder/salus/checks/adaptors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("CheckSummaryWriter", func() {
	var checkWriter CheckSummaryWriter
	var checkSummaryFilePath string

	BeforeEach(func() {
		tempFile, err := fileutil.TempFile(os.TempDir(), "checksummaryfile", "json")
		Expect(err).ToNot(HaveOccurred())

		checkSummaryFilePath = tempFile.Name()

		checkWriter = CheckSummaryWriter{
			PathToCheckSummary: checkSummaryFilePath,
			Logger:             logger.NewLogger(logger.LevelNone),
		}
	})

	AfterEach(func() {
		os.Remove(checkSummaryFilePath)
	})

	Context("Given a valid summary", func() {
		var summary adaptors.Status

		BeforeEach(func() {
			summary = adaptors.Status{
				CheckType: "check_id_1",
				CheckInfo: checks.CheckInfo{
					Status: "status",
					Note:   "note",
				},
				CheckError: "some error",
			}
		})

		It("should persist to file", func() {
			err := checkWriter.Write(summary)
			Expect(err).ToNot(HaveOccurred())

			Expect(checkSummaryFilePath).To(BeAnExistingFile())
			Expect(checkSummaryFilePath).To(BeARegularFile())

			checkSummaryFileContents, err := ioutil.ReadFile(checkSummaryFilePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(checkSummaryFileContents).ToNot(BeEmpty())

			Expect(checkSummaryFileContents).To(MatchJSON(`{
				  "CheckStatus": {
				    "check_id_1": [
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status",
					  "Note": "note"
					},
					"CheckError": "some error"
				      }
				    ]
				  }
				}
			`))
		})

		Context("Given another valid summary of the same CheckType", func() {
			var summary2 adaptors.Status

			BeforeEach(func() {
				summary2 = adaptors.Status{
					CheckType: "check_id_1",
					CheckInfo: checks.CheckInfo{
						Status: "status 2",
						Note:   "note 2",
					},
					CheckError: "some error 2",
				}
			})

			It("should persist both statuses to summary file", func() {
				err := checkWriter.Write(summary)
				Expect(err).ToNot(HaveOccurred())
				err = checkWriter.Write(summary2)
				Expect(err).ToNot(HaveOccurred())

				Expect(checkSummaryFilePath).To(BeAnExistingFile())
				Expect(checkSummaryFilePath).To(BeARegularFile())

				checkSummaryFileContents, err := ioutil.ReadFile(checkSummaryFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(checkSummaryFileContents).ToNot(BeEmpty())

				Expect(checkSummaryFileContents).To(MatchJSON(`{
				  "CheckStatus": {
				    "check_id_1": [
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status",
					  "Note": "note"
					},
					"CheckError": "some error"
				      },
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status 2",
					  "Note": "note 2"
					},
					"CheckError": "some error 2"
				      }

				    ]
				  }
				}`))
			})
		})

		Context("Given a path to a pre-existing checksummary file", func() {
			var existingCheckSummaryFile *os.File

			BeforeEach(func() {
				var err error
				existingCheckSummaryFile, err = ioutil.TempFile(os.TempDir(), "checksummarywritertestfile")
				Expect(err).ToNot(HaveOccurred())

				err = ioutil.WriteFile(existingCheckSummaryFile.Name(), []byte(`{
				  "CheckStatus": {
				    "check_id_1": [
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status",
					  "Note": "note"
					},
					"CheckError": "some error"
				      }
				    ]
				  }
				}`), 0777)
				Expect(err).ToNot(HaveOccurred())
				checkWriter = CheckSummaryWriter{
					PathToCheckSummary: existingCheckSummaryFile.Name(),
					Logger:             logger.NewLogger(logger.LevelNone),
				}
			})

			AfterEach(func() {
				os.Remove(existingCheckSummaryFile.Name())
			})

			It("should append summary contents into file", func() {
				err := checkWriter.Write(summary)
				Expect(err).ToNot(HaveOccurred())

				Expect(existingCheckSummaryFile.Name()).To(BeAnExistingFile())
				Expect(existingCheckSummaryFile.Name()).To(BeARegularFile())

				checkSummaryFileContents, err := ioutil.ReadFile(existingCheckSummaryFile.Name())
				Expect(err).ToNot(HaveOccurred())

				Expect(checkSummaryFileContents).ToNot(BeEmpty())
				Expect(checkSummaryFileContents).To(MatchJSON(`{
				  "CheckStatus": {
				    "check_id_1": [
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status",
					  "Note": "note"
					},
					"CheckError": "some error"
				      },
				      {
					"Modified": "0001-01-01T00:00:00Z",
					"CheckInfo": {
					  "Status": "status",
					  "Note": "note"
					},
					"CheckError": "some error"
				      }
				    ]
				  }
				}`))
			})

		})
	})
})
