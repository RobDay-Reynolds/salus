package writer_test

import (
	. "github.com/monkeyherder/moirai/checks/writer"

	"encoding/json"
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/golang/go/src/pkg/path/filepath"
	"github.com/golang/go/src/pkg/strconv"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/adaptors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("CheckSummaryWriter", func() {
	var checkWriter CheckSummaryWriter
	var checkSummaryFilePath string

	BeforeEach(func() {
		checkSummaryFilePath = filepath.Join(os.TempDir(), strconv.Itoa(int(time.Now().Unix())))

		checkWriter = CheckSummaryWriter{
			PathToCheckSummary: checkSummaryFilePath,
		}
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

			summaryJson, err := json.Marshal(summary)
			Expect(err).ToNot(HaveOccurred())

			Expect(checkSummaryFileContents).To(ContainSubstring(string(summaryJson)))
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

				summaryJson, err := json.Marshal(summary)
				Expect(err).ToNot(HaveOccurred())

				Expect(checkSummaryFileContents).To(ContainSubstring(string(summaryJson)))

				summaryJson2, err := json.Marshal(summary2)
				Expect(err).ToNot(HaveOccurred())

				Expect(checkSummaryFileContents).To(ContainSubstring(string(summaryJson2)))
			})
		})

		Context("Given a path to a pre-existing checksummary file", func() {
			var existingCheckSummaryFile *os.File

			BeforeEach(func() {
				var err error
				existingCheckSummaryFile, err = ioutil.TempFile(os.TempDir(), "checksummarywritertestfile")
				Expect(err).ToNot(HaveOccurred())

				checkWriter = CheckSummaryWriter{
					PathToCheckSummary: existingCheckSummaryFile.Name(),
				}

				err = checkWriter.Write(summary)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should append to file", func() {
				err := checkWriter.Write(summary)
				Expect(err).ToNot(HaveOccurred())

				Expect(existingCheckSummaryFile.Name()).To(BeAnExistingFile())
				Expect(existingCheckSummaryFile.Name()).To(BeARegularFile())

				checkSummaryFileContents, err := ioutil.ReadAll(existingCheckSummaryFile)
				Expect(err).ToNot(HaveOccurred())

				Expect(checkSummaryFileContents).ToNot(BeEmpty())
			})

		})
	})
})
