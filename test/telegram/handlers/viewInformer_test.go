package test

import (
	"algobot/internal/domain"
	"algobot/internal/domain/models"
	"algobot/internal/telegram/handlers/text"
	"algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestViewInformer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	viewFetcher := mocks2.NewMockViewFetcher(ctrl)
	serializator := mocks2.NewMockSerializator(ctrl)
	botName := "botName"

	mctx := mocks3.NewMockContext(ctrl)
	handler := text.NewViewInformer(serializator, viewFetcher, log, botName)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).AnyTimes()
	mctx.EXPECT().Message().Return(&tele.Message{Text: "/start abc"}).AnyTimes()

	t.Run("Deserialize err", func(t *testing.T) {
		errExp := errors.New("exp")
		gomock.InOrder(
			serializator.EXPECT().Deserialize("abc").Return(nil, errExp),
			mctx.EXPECT().Send("⚠️ Ошибка при расшифровке запроса!").Return(nil).Times(1),
		)

		err := handler.ServeContext(mctx)
		assert.NoError(t, err)
	})
	t.Run("Cant get action handler", func(t *testing.T) {
		errExp := errors.New("exp")
		gomock.InOrder(
			serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
				Type: 222,
				Data: []string{},
			}, errExp),
			mctx.EXPECT().Send("⚠️ Ошибка при расшифровке запроса!").Return(nil).Times(1),
		)

		err := handler.ServeContext(mctx)
		assert.NoError(t, err)
	})
	t.Run("Kids", func(t *testing.T) {
		t.Run("Happy path", func(t *testing.T) {
			gomock.InOrder(
				serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"789", "321"},
				}, nil),
				viewFetcher.EXPECT().KidView(int64(1), "789", "321").Return(models.KidView{
					Extra: "",
					Kid: models.Kid{
						FullName:   "Алексей Смирнов",
						ParentName: "Мария Смирнова",
						Email:      "alexey.smirnov@example.com",
						Phone:      "+7 (912) 345-67-89",
						Age:        10,
						BirthDate:  time.Date(2014, 3, 15, 0, 0, 0, 0, time.UTC),
						Username:   "aleksey10",
						Password:   "securepassword123",
						Groups: []models.KidViewGroup{
							{
								ID:        101,
								Title:     "Основы программирования",
								Content:   "Изучение Scratch и базовых алгоритмов",
								Status:    0,
								StartTime: time.Date(2024, 9, 1, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC),
							},
							{
								ID:        102,
								Title:     "Веб-разработка для детей",
								Content:   "HTML, CSS и основы JavaScript",
								Status:    10,
								StartTime: time.Date(2025, 1, 10, 14, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2025, 3, 30, 16, 0, 0, 0, time.UTC),
							},
							{
								ID:        103,
								Title:     "Робототехника",
								Content:   "Сборка и программирование LEGO-роботов",
								Status:    20,
								StartTime: time.Date(2025, 4, 5, 9, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2025, 6, 20, 11, 0, 0, 0, time.UTC),
							},
						},
					},
				}, nil),
				mctx.EXPECT().Send("<b>Алексей Смирнов</b>\nВозраст: 10\nДень рождения: 2014-03-15\n\n<b>Данные от аккаунта:</b>\nЛогин: <i>aleksey10</i>\nПароль: <i>securepassword123</i>\n\n<b>Родитель:</b>\nИмя: Мария Смирнова\nТелефон: +79123456789 <a href=\"https://wa.me/79123456789\">🟩 Whatsapp</a>\nПочта: alexey.smirnov@example.com\n\n<b>Группы</b>\n1 . <a href=\"https://backoffice.algoritmika.org/group/view/103\">Робототехника Сборка и программирование LEGO-роботов</a>\n🔴 Выбыл (2025-04-05 - 2025-06-20)\n\n2 . <a href=\"https://backoffice.algoritmika.org/group/view/102\">Веб-разработка для детей HTML, CSS и основы JavaScript</a>\n🟡 Переведен (2025-01-10 - 2025-03-30)\n\n3 . <a href=\"https://backoffice.algoritmika.org/group/view/101\">Основы программирования Изучение Scratch и базовых алгоритмов</a>\n🟢 Учится (2024-09-01 - 2024-12-15)\n\n", tele.ModeHTML, tele.NoPreview).Return(nil).Times(1),
			)
			err := handler.ServeContext(mctx)
			assert.NoError(t, err)
		})
		t.Run("HappyPath kid not accessebly", func(t *testing.T) {
			gomock.InOrder(
				serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"789", "321"},
				}, nil),
				viewFetcher.EXPECT().KidView(int64(1), "789", "321").Return(models.KidView{
					Extra: models.NotAccessible,
					Kid: models.Kid{
						FullName:   "Алексей Смирнов",
						ParentName: "Мария Смирнова",
						Email:      "alexey.smirnov@example.com",
						Phone:      "+7 (912) 345-67-89",
						Age:        10,
						BirthDate:  time.Date(2014, 3, 15, 0, 0, 0, 0, time.UTC),
						Username:   "aleksey10",
						Password:   "securepassword123",
						Groups: []models.KidViewGroup{
							{
								ID:        101,
								Title:     "Основы программирования",
								Content:   "Изучение Scratch и базовых алгоритмов",
								Status:    0,
								StartTime: time.Date(2024, 9, 1, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC),
							},
							{
								ID:        102,
								Title:     "Веб-разработка для детей",
								Content:   "HTML, CSS и основы JavaScript",
								Status:    10,
								StartTime: time.Date(2025, 1, 10, 14, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2025, 3, 30, 16, 0, 0, 0, time.UTC),
							},
							{
								ID:        103,
								Title:     "Робототехника",
								Content:   "Сборка и программирование LEGO-роботов",
								Status:    20,
								StartTime: time.Date(2025, 4, 5, 9, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2025, 6, 20, 11, 0, 0, 0, time.UTC),
							},
						},
					},
				}, nil),
				mctx.EXPECT().Send("⚠️ У вас больше нету доступа к ребенку\n<b>Алексей Смирнов</b>\nВозраст: 10\nДень рождения: 2014-03-15\n\n<b>Данные от аккаунта:</b>\nЛогин: <i>aleksey10</i>\nПароль: <i>securepassword123</i>\n\n<b>Родитель:</b>\nИмя: Мария Смирнова\nТелефон: +79123456789 <a href=\"https://wa.me/79123456789\">🟩 Whatsapp</a>\nПочта: alexey.smirnov@example.com\n\n<b>Группы</b>\n1 . <a href=\"https://backoffice.algoritmika.org/group/view/103\">Робототехника Сборка и программирование LEGO-роботов</a>\n🔴 Выбыл (2025-04-05 - 2025-06-20)\n\n2 . <a href=\"https://backoffice.algoritmika.org/group/view/102\">Веб-разработка для детей HTML, CSS и основы JavaScript</a>\n🟡 Переведен (2025-01-10 - 2025-03-30)\n\n3 . <a href=\"https://backoffice.algoritmika.org/group/view/101\">Основы программирования Изучение Scratch и базовых алгоритмов</a>\n🟢 Учится (2024-09-01 - 2024-12-15)\n\n", tele.ModeHTML, tele.NoPreview).Return(nil).Times(1),
			)
			err := handler.ServeContext(mctx)
			assert.NoError(t, err)
		})
		t.Run("userInfo error", func(t *testing.T) {
			t.Run("data len not 2", func(t *testing.T) {
				gomock.InOrder(
					serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
						Type: domain.UserType,
						Data: []string{"789"},
					}, nil),
					mctx.EXPECT().Send("⚠️ Невозможно получить данного ученика!").Return(nil).Times(1),
				)
				err := handler.ServeContext(mctx)
				assert.NoError(t, err)
			})
			t.Run("KidView return err", func(t *testing.T) {
				errExp := errors.New("errExp")

				gomock.InOrder(
					serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
						Type: domain.UserType,
						Data: []string{"789", "123"},
					}, nil),
					viewFetcher.EXPECT().KidView(int64(1), "789", "123").Return(models.KidView{}, errExp).Times(1),
					mctx.EXPECT().Send("⚠️ Невозможно получить данного ученика!").Return(nil).Times(1),
				)
				err := handler.ServeContext(mctx)
				assert.NoError(t, err)
			})
		})
	})
	t.Run("Groups", func(t *testing.T) {
		t.Run("Happy path", func(t *testing.T) {
			gomock.InOrder(
				serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
					Type: domain.GroupType,
					Data: []string{"123"},
				}, nil),
				viewFetcher.EXPECT().GroupView(int64(1), "123").Return(models.GroupView{
					GroupID:        1,
					GroupTitle:     "Математика для детей",
					GroupContent:   "Основы арифметики и геометрии",
					NextLessonTime: "2023-10-01T10:00:00Z",
					LessonsTotal:   12,
					LessonsPassed:  5,
					ActiveKids: []models.GroupKid{
						{
							ID:       101,
							FullName: "Иван Иванов",
							LastGroup: models.KidGroup{
								ID:        1,
								StartTime: time.Date(2023, 9, 1, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2023, 9, 1, 12, 0, 0, 0, time.UTC),
							},
						},
						{
							ID:       102,
							FullName: "Мария Петровна",
							LastGroup: models.KidGroup{
								ID:        2,
								StartTime: time.Date(2023, 9, 5, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2023, 9, 5, 12, 0, 0, 0, time.UTC),
							},
						},
						{
							ID:       103,
							FullName: "Алексей Сидоров",
							LastGroup: models.KidGroup{
								ID:        3,
								StartTime: time.Date(2023, 9, 10, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2023, 9, 10, 12, 0, 0, 0, time.UTC),
							},
						},
					},
					NotActiveKids: []models.GroupKid{
						{
							ID:       104,
							FullName: "Ольга Васильева",
							LastGroup: models.KidGroup{
								ID:        4,
								StartTime: time.Date(2023, 8, 25, 10, 0, 0, 0, time.UTC),
								EndTime:   time.Date(2023, 8, 25, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				}, nil),
				serializator.EXPECT().Serialize(domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"101", "1"},
				}).Return("1", nil).Times(1),
				serializator.EXPECT().Serialize(domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"102", "1"},
				}).Return("2", nil).Times(1),
				serializator.EXPECT().Serialize(domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"103", "1"},
				}).Return("3", nil).Times(1),
				serializator.EXPECT().Serialize(domain.SerializeMessage{
					Type: domain.UserType,
					Data: []string{"104", "1"},
				}).Return("4", nil).Times(1),

				mctx.EXPECT().Send("<a href=\"https://backoffice.algoritmika.org/group/view/1\">Математика для детей Основы арифметики и геометрии</a>\n\n<b>Следующая лекция</b>: 2023-10-01T10:00:00Z\n<b>Всего пройдено</b> 5 лекций из 12\n\nАктивные дети: 3 | Выбыло: 1 | Всего: 4\n<b>Активные дети</b>:\n1. <a href=\"https://t.me/botName?start=1\">Иван Иванов</a>\n2. <a href=\"https://t.me/botName?start=2\">Мария Петровна</a>\n3. <a href=\"https://t.me/botName?start=3\">Алексей Сидоров</a>\n<b>Выбыли дети</b>:\n1. <a href=\"https://t.me/botName?start=4\">Ольга Васильева</a> (🟡 Переведен: 2023-08-25)\n", tele.ModeHTML, tele.NoPreview).Return(nil).Times(1),
			)
			err := handler.ServeContext(mctx)
			assert.NoError(t, err)
		})
		t.Run("groupInfo error", func(t *testing.T) {
			t.Run("GroupView return err", func(t *testing.T) {
				errExp := errors.New("errExp")

				gomock.InOrder(
					serializator.EXPECT().Deserialize("abc").Return(&domain.SerializeMessage{
						Type: domain.GroupType,
						Data: []string{"123"},
					}, nil),
					viewFetcher.EXPECT().GroupView(int64(1), "123").Return(models.GroupView{}, errExp).Times(1),
					mctx.EXPECT().Send("⚠️ Невозможно получить данную группу!").Return(nil).Times(1),
				)
				err := handler.ServeContext(mctx)
				assert.NoError(t, err)
			})
		})
	})

}
