package shift_test

import (
	"context"
	"math"
	"testing"
	"time"
	"williamfeng323/mooncake-duty/src/domains/account"
	"williamfeng323/mooncake-duty/src/domains/project"
	"williamfeng323/mooncake-duty/src/domains/shift"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

func TestShift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shift domain Suite")
}

var _ = Describe("Shift", func() {
	prj := project.NewProject("TestShift", "Project for testing the shift")
	prj.Create()

	AfterSuite(func() {
		repoimpl.GetProjectRepo().DeleteOne(context.Background(), bson.M{"_id": prj.ID})
	})
	Describe("#NewShift", func() {
		Describe("If project does not exist", func() {
			It("should return project not found error and nil shift", func() {
				shift, err := shift.NewShift(primitive.NewObjectID(), shift.Mon, time.Now(), shift.Weekly)
				Expect(shift).To(BeNil())
				Expect(err).To(MatchError(project.NotFoundError{}))
			})
		})
		Describe("If the project exists", func() {
			Context("When the start date is not the first date of the shift", func() {
				Context("When the recurrence is daily", func() {
					It("should replace it to the first date of the shift base on weekStart", func() {
						shiftStartedDate := time.Date(2020, 6, 16, 0, 0, 0, 0, time.Now().Location())
						shift, err := shift.NewShift(prj.ID, shift.Mon, shiftStartedDate, shift.Daily)
						Expect(shift).ToNot(BeNil())
						Expect(err).To(BeNil())
						Expect(shift.ShiftFirstDate.Format(time.RFC3339)).To(Equal(utils.ToDateStarted(shiftStartedDate).Format(time.RFC3339)))
					})
				})
				Context("When the recurrence is weekly", func() {
					It("should replace it to the first date of the shift base on weekStart", func() {
						shift, err := shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 6, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
						Expect(shift).ToNot(BeNil())
						Expect(err).To(BeNil())
						Expect(shift.ShiftFirstDate.Format(time.RFC3339)).To(Equal("2020-06-15T00:00:00+08:00"))
					})
				})
				Context("When the recurrence is bi-weekly", func() {
					It("should replace it to the first date of the shift base on weekStart", func() {
						shift, err := shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 6, 16, 0, 0, 0, 0, time.Now().Location()), shift.BiWeekly)
						Expect(shift).ToNot(BeNil())
						Expect(err).To(BeNil())
						Expect(shift.ShiftFirstDate.Format(time.RFC3339)).To(Equal("2020-06-15T00:00:00+08:00"))
					})
				})
			})
		})
	})
	Describe("#DefaultShifts", func() {
		var testShift *shift.Shift

		Context("When no member in the domain", func() {
			It("should return nil shift and no member error", func() {
				testShift, _ = shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 6, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
				tempShift, err := testShift.DefaultShifts(time.Now(), time.Now())
				Expect(tempShift).To(BeNil())
				Expect(err).To(MatchError(shift.NoMemberError{}))
			})
		})
		Describe("Giving members configured properly in the domain", func() {
			acct1, _ := account.NewAccount("Test1@test.com", "Testaccount1")
			acct2, _ := account.NewAccount("Test2@test.com", "Testaccount2")
			acct3, _ := account.NewAccount("Test3@test.com", "Testaccount3")
			acct4, _ := account.NewAccount("Test4@test.com", "Testaccount4")
			acct5, _ := account.NewAccount("Test5@test.com", "Testaccount5")
			acct6, _ := account.NewAccount("Test6@test.com", "Testaccount6")
			acct7, _ := account.NewAccount("Test7@test.com", "Testaccount7")
			acctSet := []*account.Account{acct1, acct2, acct3, acct4, acct5, acct6, acct7}

			BeforeEach(func() {
				for _, v := range acctSet {
					v.Save(false)
				}
				testShift, _ = shift.NewShift(prj.ID, shift.Mon, time.Date(2020, 5, 16, 0, 0, 0, 0, time.Now().Location()), shift.Weekly)
				testShift.T1Members = []primitive.ObjectID{acct1.ID, acct2.ID, acct3.ID}
				testShift.T2Members = []primitive.ObjectID{acct4.ID, acct5.ID, acct6.ID}
				testShift.T3Members = []primitive.ObjectID{acct7.ID}
			})
			AfterEach(func() {
				acctRepo := repoimpl.GetAccountRepo()
				for _, v := range acctSet {
					acctRepo.DeleteOne(context.Background(), bson.M{"_id": v.ID})
				}
			})
			Context("When the period start date before the shift first date", func() {
				It("should return period out of scope error", func() {
					tempShift, err := testShift.DefaultShifts(testShift.ShiftFirstDate.AddDate(0, 0, -4), time.Now())
					Expect(tempShift).To(BeNil())
					Expect(err).To(MatchError(shift.OutOfScopeError{}))
				})
			})
			Context("When the period start date after the period end date", func() {
				It("should return period out of scope error", func() {
					tempShift, err := testShift.DefaultShifts(time.Now().AddDate(0, 0, 4), time.Now())
					Expect(tempShift).To(BeNil())
					Expect(err).To(MatchError(shift.OutOfScopeError{}))
				})
			})
			Context("When the recurrence is daily", func() {
				It("should return list of shifts for every day within the period", func() {
					testShift.ShiftRecurrence = shift.Daily
					tempShift, err := testShift.DefaultShifts(utils.ToDateStarted(time.Now().AddDate(0, 0, -10)), time.Now())
					Expect(err).To(BeNil())
					Expect(len(tempShift)).To(Equal(10))
				})
				It("should return list of shifts calculated base on the started date of shift", func() {
					testShift.ShiftRecurrence = shift.Daily
					testShift.ShiftFirstDate = utils.ToDateStarted(time.Now().AddDate(0, 0, -5))
					tempShifts, err := testShift.DefaultShifts(utils.ToDateStarted(time.Now().AddDate(0, 0, -4)), time.Now())
					Expect(err).To(BeNil())
					Expect(len(tempShifts)).To(Equal(4))
					Expect(tempShifts[0].T1Member.String()).To(Equal(acct2.ID.String()))
					Expect(tempShifts[0].T2Member.String()).To(Equal(acct5.ID.String()))
					Expect(tempShifts[0].T3Member.String()).To(Equal(acct7.ID.String()))

					Expect(tempShifts[1].T1Member.String()).To(Equal(acct3.ID.String()))
					Expect(tempShifts[1].T2Member.String()).To(Equal(acct6.ID.String()))
					Expect(tempShifts[1].T3Member.String()).To(Equal(acct7.ID.String()))

					Expect(tempShifts[2].T1Member.String()).To(Equal(acct1.ID.String()))
					Expect(tempShifts[2].T2Member.String()).To(Equal(acct4.ID.String()))
					Expect(tempShifts[2].T3Member.String()).To(Equal(acct7.ID.String()))

					Expect(tempShifts[3].T1Member.String()).To(Equal(acct2.ID.String()))
					Expect(tempShifts[3].T2Member.String()).To(Equal(acct5.ID.String()))
					Expect(tempShifts[3].T3Member.String()).To(Equal(acct7.ID.String()))
				})
			})
			Context("When the recurrence is weekly", func() {
				BeforeEach(func() {
					testShift.ShiftFirstDate = utils.FirstDateOfWeek(time.Now().AddDate(0, 0, -28), time.Weekday(shift.Mon))
					testShift.WeekStart = shift.Mon
					testShift.ShiftRecurrence = shift.Weekly
				})
				Context("giving the week start on Mon", func() {
					It("The start date of the shift should be Mon and the end date should be Sun", func() {
						tempShifts, err := testShift.DefaultShifts(utils.ToDateStarted(time.Now().AddDate(0, 0, -21)), time.Now())
						Expect(err).To(BeNil())
						Expect(len(tempShifts)).To(Equal(4))
						for _, v := range tempShifts {
							Expect(v.StartDate.Weekday()).To(Equal(time.Monday))
							Expect(v.EndDate.Weekday()).To(Equal(time.Sunday))
						}
					})
				})
				Context("giving the week start on Sun", func() {
					It("The start date of the shift should be Sun and the end date should be Sat", func() {
						testShift.WeekStart = shift.Sun
						testShift.ShiftFirstDate = utils.FirstDateOfWeek(time.Now().AddDate(0, 0, -28), time.Sunday)
						tempShifts, err := testShift.DefaultShifts(utils.ToDateStarted(time.Now().AddDate(0, 0, -21)), time.Now())
						Expect(err).To(BeNil())
						Expect(len(tempShifts)).To(Equal(4))
						for _, v := range tempShifts {
							Expect(v.StartDate.Weekday()).To(Equal(time.Sunday))
							Expect(v.EndDate.Weekday()).To(Equal(time.Saturday))
							Expect(math.Ceil(v.EndDate.Sub(v.StartDate).Hours() / 24)).To(Equal(float64(7)))
						}
					})
				})
				Context("giving the shift start date 4 weeks before and the requesting start date of the period is 3 weeks before", func() {
					It("should assign the shift to members by week and shifts should rotate from the starting week", func() {
						tempShifts, err := testShift.DefaultShifts(utils.ToDateStarted(time.Now().AddDate(0, 0, -21)), time.Now())
						Expect(len(tempShifts)).To(Equal(4))
						Expect(err).To(BeNil())
						Expect(tempShifts[0].T1Member.String()).To(Equal(acct2.ID.String()))
						Expect(tempShifts[0].T2Member.String()).To(Equal(acct5.ID.String()))
						Expect(tempShifts[0].T3Member.String()).To(Equal(acct7.ID.String()))

						Expect(tempShifts[1].T1Member.String()).To(Equal(acct3.ID.String()))
						Expect(tempShifts[1].T2Member.String()).To(Equal(acct6.ID.String()))
						Expect(tempShifts[1].T3Member.String()).To(Equal(acct7.ID.String()))

						Expect(tempShifts[2].T1Member.String()).To(Equal(acct1.ID.String()))
						Expect(tempShifts[2].T2Member.String()).To(Equal(acct4.ID.String()))
						Expect(tempShifts[2].T3Member.String()).To(Equal(acct7.ID.String()))
					})
				})
			})
		})
	})
})
