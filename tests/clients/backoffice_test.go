package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"tgbot/internal/clients"
	"time"
)

func TestBackoffice(t *testing.T) {
	boSettings := clients.BackofficeSetting{
		Retry:        3,
		Timeout:      50 * time.Millisecond,
		RetryTimeout: 50 * time.Millisecond,
	}

	t.Run("GetKidsNamesByGroup", func(t *testing.T) {
		t.Run("GENERAL | 401 | Unauthorized", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("[]"))
			}))

			bo := clients.NewBackoffice(ts.URL, boSettings)
			_, err := bo.GetKidsNamesByGroup("", "")
			assertError(t, err, 401, "[]")
		})
		t.Run("GENERAL | 500 | Servers returns error", func(t *testing.T) {
			var calls []string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				calls = append(calls, r.Method)
			}))
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			_, err := bo.GetKidsNamesByGroup("", "")
			assertError(t, err, 500, "")
			if len(calls) != 3 {
				t.Fatalf("expected 3 calls, got %d", len(calls))
			}

		})
		t.Run("GENERAL | Timeout | Servers return timeout", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			_, err := bo.GetKidsNamesByGroup("", "")
			assertError(t, err, 500, "_")
		})
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			groupId := "333"

			ts := getBOServer(t, map[string]string{
				"groupId": groupId,
				"expand":  "lastGroup",
			}, "GET", cookie, "", GetKidsNamesByGroupResponse)
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			kids, err := bo.GetKidsNamesByGroup(cookie, groupId)
			assertNoError(t, err)

			ks := clients.GroupResponse{}
			json.Unmarshal([]byte(GetKidsNamesByGroupResponse), &ks)
			if reflect.DeepEqual(ks, kids) {
				t.Errorf("Wanted: %+v", ks)
				t.Errorf("Got: %+v", kids)
				t.Fatalf("expected kids to be different")
			}
		})
	})
	t.Run("GetKidsStatsByGroup", func(t *testing.T) {
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			groupId := "333"

			ts := getBOServer(t, map[string]string{
				"group": groupId,
			}, "GET", cookie, "", GetKidsStatsByGroupResponse)
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			kids, err := bo.GetKidsStatsByGroup(cookie, groupId)
			assertNoError(t, err)

			ks := clients.KidsStats{}
			json.Unmarshal([]byte(GetKidsStatsByGroupResponse), &ks)
			if reflect.DeepEqual(ks, kids) {
				t.Errorf("Wanted: %+v", ks)
				t.Errorf("Got: %+v", kids)
				t.Fatalf("expected kids to be different")
			}
		})
	})
	t.Run("GetKidsMessages", func(t *testing.T) {
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			ts := getBOServer(t, map[string]string{
				"from":  "0",
				"limit": "30",
			}, "GET", cookie, "", KidsMessagesResponse)
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			kids, err := bo.GetKidsMessages(cookie)
			assertNoError(t, err)

			ks := clients.KidsStats{}
			json.Unmarshal([]byte(KidsMessagesResponse), &ks)
			if reflect.DeepEqual(ks, kids) {
				t.Errorf("Wanted: %+v", ks)
				t.Errorf("Got: %+v", kids)
				t.Fatalf("expected kids to be different")
			}
		})
	})
	t.Run("GetAllGroupsByUser", func(t *testing.T) {
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			ts := getBOServer(t, map[string]string{
				"GroupSearch[status][]": "active",
				"presetType":            "all",
				"_pjax":                 "#group-grid-pjax",
			}, "GET", cookie, "", HtmlResponse)
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			kids, err := bo.GetAllGroupsByUser(cookie)
			assertNoError(t, err)

			if !reflect.DeepEqual(Kids, kids) {
				t.Errorf("Wanted: %+v", Kids)
				t.Fatalf("Got: %+v", kids)
			}
		})
	})
	t.Run("CloseLession", func(t *testing.T) {
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			group := "111"
			lession := "222"
			ts := getBOServer(t, map[string]string{}, "POST", cookie, "ajaxUrl=^%^2Fapi^%^2Fv2^%^2Fgroup^%^2Flesson^%^2Fstatus&btnClass=btn+btn-xs+btn-danger&status=0&lessonId=222&groupId=111", "[]")
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			err := bo.CloseLession(cookie, group, lession)
			assertNoError(t, err)
		})
	})
	t.Run("OpenLession", func(t *testing.T) {
		t.Run("200 | Servers return OK", func(t *testing.T) {
			cookie := "111"
			group := "111"
			lession := "222"
			ts := getBOServer(t, map[string]string{}, "POST", cookie, "ajaxUrl=^%^2Fapi^%^2Fv2^%^2Fgroup^%^2Flesson^%^2Fstatus&btnClass=btn+btn-xs+btn-danger&status=10&lessonId=222&groupId=111", "[]")
			defer ts.Close()

			bo := clients.NewBackoffice(ts.URL, boSettings)
			err := bo.OpenLession(cookie, group, lession)
			assertNoError(t, err)
		})
	})
}

func getBOServer(t *testing.T, wantedParams map[string]string, wantedMethod string, wantedCookie string, wantedBody string, serverResponse string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri, _ := url.Parse(r.RequestURI)
		params := uri.Query()

		for k, v := range wantedParams {
			if params.Get(k) != v {
				t.Fatalf("expected %s=%s, got %s", k, v, params.Get(k))
			}
		}
		if r.Method != wantedMethod {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("Cookie") != wantedCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if getString(r.Body) != wantedBody {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverResponse))
	}))
	return ts
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func assertError(t *testing.T, err error, i int, s string) {
	t.Helper()

	// Приводим ошибку к типу *clients.ClientError
	clientErr, ok := err.(*clients.ClientError)
	if !ok {
		t.Errorf("Expected error of type *clients.ClientError, but got: %T", err)
		return
	}

	if clientErr == nil {
		t.Fatal("Expected error, got nil")
	}
	if clientErr.Code != i {
		t.Errorf("%+v\n", err)
		t.Fatalf("Expected code %d, got %d", i, clientErr.Code)
	}
	if s != "_" {
		if clientErr.Message != s {
			t.Errorf("%+v\n", err)
			t.Fatalf("Expected message %s, got %s", s, clientErr.Message)
		}
	}
}
func getString(body io.Reader) string {
	all, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(all)
}

const GetKidsNamesByGroupResponse = `{"status": "success","data": {"items": [{"id": 70245813,"firstName": "firstName","lastName": "lastName","fullName": "firstName lastName","parentName": "parentName","email": "email@mail.ru","hasLaptop": -1,"phone": "+7 (999) 999-99-99","age": 10,"birthDate": "2014-12-16T00:00:00+03:00","createdAt": "2024-11-18T08:18:45+00:00","updatedAt": "2024-12-03T08:31:26+00:00","deletedAt": null,"hasBranchAccess": true,"username": "username","password": "password","lastGroup": {"id": 98637162,"groupStudentId": 6553709,"title": "title","content": "content","track": 2,"status": 0,"startTime": "2024-11-25T10:54:55+03:00","endTime": "9999-12-31T00:00:00+03:00","courseId": 729,"createdAt": "2024-09-17T07:41:22+00:00","updatedAt": "2025-01-15T09:10:46+00:00","deletedAt": null},"_links": {"self": {"href": "/student/update/70245813"}}}]}}`
const GetKidsStatsByGroupResponse = `{"status": "success","data": [{"student_id": 3356212,"attendance": [{"lesson_id": 4560,"lesson_title": "1","start_time_formatted": "вс 22.09.24 14:00","status": "absent"},{"lesson_id": 4561,"lesson_title": "2","start_time_formatted": "вс 29.09.24 14:00","status": "absent"}]}]}`
const KidsMessagesResponse = `{"status": "success","data": {"projects": [{"uid": "1039632level154551","new": false,"senderId": 1039632,"senderScope": "student","type": "text","content": "11","name": "a a","lastTime": "11 янв. 15:25","title": "Что такое w-w w?","link": "/task-preview/16402?student=s&lesson=s&position=5&s=1&groupId=98619873"}]}}`
const HtmlResponse = `<!DOCTYPE html><html lang="ru-RU" dir="ltr"><head><meta charset="UTF-8" /><meta name="viewport" content="width=device-width, initial-scale=1" /><title>Группы</title><link type="image/x-icon" href="/favicon.ico?v=3" rel="icon"><link href="/assets/b3de4d09/themes/smoothness/jquery-ui.css?v=1737106652" rel="stylesheet"><link href="/assets/eedc24e4/css/font-awesome.min.css?v=1737106652" rel="stylesheet"><link href="/assets/9793f2e6/css/bootstrap.css?v=1737106652" rel="stylesheet"><link href="/assets/45d70904/css/comments.css?v=1737106652" rel="stylesheet"><script src="/assets/405e6a63/js/dialog.min.js?v=1737106652"></script><!-- End Google Tag Manager --></head><body class="hold-transition skin-dark-blue sidebar-mini sidebar-collapse"><div class="wrapper"><header class="defaultHandler-header" id="header"><!-- Logo --><a href="/" class="logo"><!-- mini logo for sidebar mini 50x50 pixels --><span class="logo-mini"><imgsrc="/img/logo_v2/rgb-logo-znak-p.png?v=2"alt="Алгоритмика"/></span><!-- logo for regular state and mobile devices --><span class="logo-lg"><imgsrc="/img/logo_v2/rgb-logo-small-hor-p.png?v=2"alt="Алгоритмика"/></span></a><!-- Sidebar toggle button--><a href="#" class="sidebar-toggle" data-toggle="push-menu" role="button"><span class="sr-only">Toggle navigation</span><span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span></a><!-- Header Navbar: style can be found in header.less --><nav class="navbar navbar-static-top"><search-widget></search-widget><div class="navbar-custom-menu"><ul class="nav navbar-nav"><!-- Dev Widget --><!-- Closest lesson widget --><li class="dropdown closest-lesson"><a href="#" class="dropdown-toggle" data-toggle="dropdown"><i class="fa fa-bullhorn"></i><span class="hidden-xs">Ближайший урок</span></a><ul class="dropdown-menu closest-lesson__dropdown"><a href="/group/view/98619913"><div class="closest-lesson-widget"><h4 class="closest-lesson-widget__title">Ближайший урок:</h4><p>18.01.2025 10:00 (суббота)</p><p>Библиотека № 7 сб 10.00</p><p>Группа по курсу КГ</p><p>Структурируем презентацию</p></div></a></ul></li><!-- Branches --><li class="branch-menu flip new-ui"><show-branch-scope-selector-button></show-branch-scope-selector-button></li><!-- User Account: style can be found in dropdown.less --><li class="dropdown user user-menu"><a href="#" class="dropdown-toggle" data-toggle="dropdown"><i class="fa fa-user"></i><span class="hidden-xs">Данил Павлов</span></a><ul class="dropdown-menu"><li><a href="/user/profile">Профиль</a>                </li><li><li class="dropdown language-menu"><a href="#" class="js-language-menu-button"><i class="fa fa-flag"></i><span class="hidden-xs">Русский</span></a><ul class="dropdown-menu"><li><span class="dropdown-item is-selected">Русский</span></li><li><aclass="js-switch-language"href="#"data-locale="az-AZ">Azərbaycan            </a></li><li><aclass="js-switch-language"href="#"data-locale="en-US">English            </a></li><li><aclass="js-switch-language"href="#"data-locale="he-HE">עִבְרִית (Иврит)            </a></li><li><aclass="js-switch-language"href="#"data-locale="tt-RU">Tatar            </a></li><li><aclass="js-switch-language"href="#"data-locale="es-419">Español            </a></li><li><aclass="js-switch-language"href="#"data-locale="el-GR">ελληνικά (Греческий)            </a></li><li><aclass="js-switch-language"href="#"data-locale="ar-SA">العربية (арабский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="pt-BR">Português            </a></li><li><aclass="js-switch-language"href="#"data-locale="zh-CN">中文 (китайский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="uk-UA">Українська            </a></li><li><aclass="js-switch-language"href="#"data-locale="de-DE">Deutsch            </a></li><li><aclass="js-switch-language"href="#"data-locale="kk-KZ">Қазақша            </a></li><li><aclass="js-switch-language"href="#"data-locale="ka-GE">ქართული (Грузинский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="cs-CZ">Česky            </a></li><li><aclass="js-switch-language"href="#"data-locale="sk">Slovak (slovenčina)            </a></li><li><aclass="js-switch-language"href="#"data-locale="fr">French            </a></li><li><aclass="js-switch-language"href="#"data-locale="tr-TR">Türkçe            </a></li><li><aclass="js-switch-language"href="#"data-locale="af-AF">Afrikaans            </a></li><li><aclass="js-switch-language"href="#"data-locale="tr">Tr            </a></li><li><aclass="js-switch-language"href="#"data-locale="tg-TG">забо́ни тоҷикӣ (Таджикский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="ro-RO">Română            </a></li><li><aclass="js-switch-language"href="#"data-locale="pl-PL">Polski            </a></li><li><aclass="js-switch-language"href="#"data-locale="it-IT">italiano            </a></li><li><aclass="js-switch-language"href="#"data-locale="da-DA">dansk (Danish)            </a></li><li><aclass="js-switch-language"href="#"data-locale="et-EE">Eesti keel (эстонский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="ca-CA">Català (каталанский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="cr-CR">Crnogorski jezik (черногорский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="lt-LT">lietuvių kalba (литовский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="mn-MN">mongɣol xele            </a></li><li><aclass="js-switch-language"href="#"data-locale="nl-NL">de Nederlandse taal            </a></li><li><aclass="js-switch-language"href="#"data-locale="bg-BG">български език            </a></li><li><aclass="js-switch-language"href="#"data-locale="hr-HR">hrvatski jezik            </a></li><li><aclass="js-switch-language"href="#"data-locale="bs-BS">bosanski jezik            </a></li><li><aclass="js-switch-language"href="#"data-locale="sr-sr">српски језик (сербский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="th-TH">ภาษาไทย (тайский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="sq-sq">Gjuha shqipe (албанский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="mk-mk">македонски јазик            </a></li><li><aclass="js-switch-language"href="#"data-locale="hu-HU">magyar (венгерский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="id-ID">Bahasa Indonesia (индонезийский)            </a></li><li><aclass="js-switch-language"href="#"data-locale="kh-KH">ភាសាខ្មែរ  ខេមរភាសា (Кхмерский язык)            </a></li><li><aclass="js-switch-language"href="#"data-locale="vi-VN">Vietnamese            </a></li></ul></li></li><li><a href="/s/auth/api/e/user/logout">Выйти</a></li></ul></li></ul></div></nav></header><!-- Left side column. contains the logo and sidebar --><aside class="defaultHandler-sidebar"><!-- sidebar: style can be found in sidebar.less --><section class="sidebar"><ul class="sidebar-menu" data-widget="tree"><li><a href="/group/default/schedule"><i class="fa fa-clock-o"></i> <span>Расписание</span></a></li><li><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;GroupSearch%5Bstatus%5D%5B1%5D=recruiting&amp;presetType=all"><i class="fa fa-graduation-cap"></i> <span>Группы</span></a></li><li><a href="/student"><i class="fa fa-child"></i> <span>Ученики</span></a></li><li><a href="/venue"><i class="fa fa-bank"></i> <span>Площадки</span></a></li><li><a href="/course"><i class="fa fa-archive"></i> <span>Курсы</span></a></li><li id="knowledgebase-menu"></li></ul>        </section><!-- /.sidebar --></aside><div class="content-wrapper"><div class="container-fluid"><div id="app-container" data-user-id="42407"><portal-target name="dialog" multiple /></div><div class="row form-group"></div><div class="row"><div class="col-lg-12"><div class="group-index"><div class="alg-page__header"><h3>Группы</h3><div class="group-create__actions"><div id="group-create-button"></div></div></div><div class="panel panel-default"><div class="panel-body"><div id="group-grid-pjax" data-pjax-container="" data-pjax-push-state><div class="kv-loader-overlay"><div class="kv-loader"></div></div><div id="group-grid" class="grid-view is-bs3 kv-grid-bs3 hide-resize" data-krajee-grid="kvGridInit_33dcf33a" data-krajee-ps="ps_group_grid_container"><div class="row quick-actions"><div class="col-sm-6 text-left flip"><div class="quick-actions "><a class="filter-link is-active js-filter-link" href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;GroupSearch%5Bstatus%5D%5B1%5D=recruiting&amp;presetType=all" data-counter-url="/api/v2/group/preset/count">Все (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5Btype%5D%5B0%5D=masterclass&amp;GroupSearch%5Btype%5D%5B1%5D=regular&amp;GroupSearch%5Btype%5D%5B2%5D=intensive&amp;GroupSearch%5Btype%5D%5B3%5D=individual&amp;GroupSearch%5BcourseContentType%5D%5B0%5D=course&amp;GroupSearch%5BcourseContentType%5D%5B1%5D=intensive&amp;GroupSearch%5BcourseStartTime%5D=2024-01-17+-+2025-01-17&amp;GroupSearch%5BcourseEndTime%5D=2025-01-17+-+2026-01-17&amp;presetType=active" data-counter-url="/api/v2/group/preset/count">Активные (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5Btype%5D%5B0%5D=masterclass&amp;GroupSearch%5Btype%5D%5B1%5D=regular&amp;GroupSearch%5Btype%5D%5B2%5D=intensive&amp;GroupSearch%5Btype%5D%5B3%5D=individual&amp;GroupSearch%5BcourseContentType%5D%5B0%5D=course&amp;GroupSearch%5BcourseContentType%5D%5B1%5D=intensive&amp;GroupSearch%5BcourseStartTime%5D=2025-01-17+-+2026-01-17&amp;sort=priorityAndNextLessonTime&amp;presetType=recruitmentInGroup" data-counter-url="/api/v2/group/preset/count">Набор в группы (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5Btype%5D%5B0%5D=masterclass&amp;GroupSearch%5Btype%5D%5B1%5D=regular&amp;GroupSearch%5Btype%5D%5B2%5D=intensive&amp;GroupSearch%5Btype%5D%5B3%5D=individual&amp;GroupSearch%5BcourseContentType%5D=masterclass&amp;GroupSearch%5BcourseStartTime%5D=2025-01-17+-+2026-01-17&amp;sort=priorityAndNextLessonTime&amp;presetType=recruitmentInMasterClass" data-counter-url="/api/v2/group/preset/count">Набор в МК (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5Btype%5D%5B0%5D=masterclass&amp;GroupSearch%5Btype%5D%5B1%5D=regular&amp;GroupSearch%5Btype%5D%5B2%5D=intensive&amp;GroupSearch%5Btype%5D%5B3%5D=individual&amp;GroupSearch%5BcourseContentType%5D%5B0%5D=course&amp;GroupSearch%5BcourseContentType%5D%5B1%5D=intensive&amp;GroupSearch%5BcourseStartTime%5D=2025-01-17+-+2026-01-17&amp;GroupSearch%5Bis_online%5D=1&amp;sort=priorityAndNextLessonTime&amp;presetType=recruitmentInOnlineGroup" data-counter-url="/api/v2/group/preset/count">Набор в онлайн группы (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5Btype%5D%5B0%5D=masterclass&amp;GroupSearch%5Btype%5D%5B1%5D=regular&amp;GroupSearch%5Btype%5D%5B2%5D=intensive&amp;GroupSearch%5Btype%5D%5B3%5D=individual&amp;GroupSearch%5BcourseContentType%5D=masterclass&amp;GroupSearch%5BcourseStartTime%5D=2025-01-17+-+2026-01-17&amp;GroupSearch%5Bis_online%5D=online&amp;sort=priorityAndNextLessonTime&amp;presetType=recruitmentInOnlineMasterClass" data-counter-url="/api/v2/group/preset/count">Набор в онлайн МК (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a><span class="note">|</span><a class="filter-link js-filter-link" href="/group?GroupSearch%5BpresetType%5D=deleted&amp;presetType=deleted" data-counter-url="/api/v2/group/preset/count">Удаленные (<span class="filter-link__counter ltr-text" dir="ltr">-</span>)</a></div><br/><span class="ltr-text" dir="ltr">Показаны записи <b>1-10</b> из <b>10</b></span></div><div class="col-sm-6 text-right flip"><div class="group-grid-column-picker btn-group"><buttontype="button"class="dropdown-toggle btn btn-sm btn-default ui-button ui-state-default ui-corner-all ui-button-text-only"data-toggle="dropdown"aria-haspopup="true"aria-expanded="false"><span class="ui-button-text"><span class="btn-text" id="lesson-status-span">Колонки</span><span class="caret"></span></span></button><form class="dropdown dropdown-menu" id="group-grid-column-picker" style="width:15em; padding:0.5em;"><div><label><input type="checkbox" name="group-grid-column-picker[]" value="id" checked> ID</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="title" checked> Название</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="venue" checked> Площадка</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="venueLocation"> География</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="activeStudentCount" checked> Уч-ки</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="paidStudentCount" checked>  Оплатившие</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="transferredStudentCount" checked>  Зачисленные</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="expelledStudentCount">  Отчисленные</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="paidStudentAfterMkCount"> AMC Ученики оплатившие после мастеркласса</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="expelledAndTransferredStudentCount">  Отчисленные и переведенные</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="firstLessonTime"> Первый урок</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="nextLessonTime" checked> Время след. урока</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="nextLessonNumber" checked>  След. урок</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="nextLessonTitle" checked> След. урок</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="teacher" checked> Преподаватель</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="curator" checked> Куратор</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="tutor"> Тьютор</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="priority_level"> Приоритет</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="type" checked> Тип группы</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="status" checked> Статус</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="course"> Курс</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="courseType"> Тип курса</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="courseContentType"> Тип контента</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="courseStartTime"> Время начала курса</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="courseEndTime"> Время окончания курса</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="is_full"> Доступность</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="is_online" checked> Формат</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="created_at"> Создан</label><br><label><input type="checkbox" name="group-grid-column-picker[]" value="created_by"> Кем создан</label></div>        <button class="btn btn-xs btn-primary" type="submit">Сохранить</button><a href="#" class="btn btn-xs btn-default cancel">Сбросить</a></form></div><div class="form-inline pull-right flip"><span style="display: none" id="group-grid-clear-filters-btn"class="btn btn-sm btn-default">Очистить фильтры        </span>Записей на странице<select class="form-control input-sm" name="grid-page-size"><option value="20" selected>20</option><option value="100">100</option><option value="200">200</option></select></div></div></div><div id="group-grid-container" class="kv-grid-container"><table class="kv-grid-table table table-hover" style="background: white"><colgroup><col style="width:1%"><col style="width:15%"><col style="width:20%"><col><col><col style="width:10%"><col style="width:1%; text-align:center"><col style="width:15%"><col style="width:15%"><col style="width:15%"><col><col><col></colgroup><thead class="kv-table-header group-grid"><tr><th data-col-seq="id"><a class="kv-sort-link desc" href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=id" data-sort="id">ID<span class="kv-sort-icon"><i class="glyphicon glyphicon-sort-by-attributes-alt"></i></span></a></th><th data-col-seq="title"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=title" data-sort="title">Название</a></th><th data-col-seq="venue">Площадка</th><th title="Учеников в группе / с ноутбуками" data-col-seq="activeStudentCount"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=active_student_count" data-sort="active_student_count">Уч-ки</a></th><th title="Ученики зачисленные в регулярную группу" data-col-seq="transferredStudentCount"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=transferredStudentCount" data-sort="transferredStudentCount"><span class="fa fa-graduation-cap"></span><span class="label-collapsed"> Зачисленные</span></a></th><th data-col-seq="nextLessonTime"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=nextLessonTime" data-sort="nextLessonTime">Время след. урока</a></th><th title="След. урок #" data-col-seq="nextLessonNumber"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=nextLessonNumber" data-sort="nextLessonNumber"><span class="fa fa-hashtag"></span><span class="label-collapsed"> След. урок</span></a></th><th data-col-seq="nextLessonTitle"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=nextLessonTitle" data-sort="nextLessonTitle">След. урок</a></th><th data-col-seq="teacher">Преподаватель</th><th data-col-seq="curator"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=curator" data-sort="curator">Куратор</a></th><th data-col-seq="type"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=type" data-sort="type">Тип группы</a></th><th data-col-seq="status"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=status" data-sort="status">Статус</a></th><th data-col-seq="is_online"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;sort=is_online" data-sort="is_online">Формат</a></th></tr><tr id="group-grid-filters" class="filters skip-export"><td data-col-seq="id"><input type="text" class="form-control" name="GroupSearch[id]"></td><td data-col-seq="title"><input type="text" class="form-control" name="GroupSearch[title]"></td><td data-col-seq="venue"><input type="text" class="form-control" name="GroupSearch[venue]"></td><td data-col-seq="activeStudentCount"><input type="text" class="form-control" name="GroupSearch[active_student_count]"></td><td data-col-seq="transferredStudentCount">&nbsp;</td><td data-col-seq="nextLessonTime"><input type="text" id="groupsearch-nextlessontime" class="form-control" name="GroupSearch[nextLessonTime]" value="" data-krajee-daterangepicker="daterangepicker_eb5279fd"></td><td data-col-seq="nextLessonNumber"><input type="text" class="form-control" name="GroupSearch[nextLessonNumber]"></td><td data-col-seq="nextLessonTitle"><input type="text" class="form-control" name="GroupSearch[nextLessonTitle]"></td><td data-col-seq="teacher"><input type="text" class="form-control" name="GroupSearch[teacher]"></td><td data-col-seq="curator"><input type="text" class="form-control" name="GroupSearch[curator]"></td><td data-col-seq="type"><div class="kv-plugin-loading loading-groupsearch-type">&nbsp;</div><input type="hidden" name="GroupSearch[type]" value=""><select id="groupsearch-type" class="form-control" name="GroupSearch[type][]" multiple size="4" data-s2-options="s2options_3feda8aa" data-krajee-select2="select2_08b55e27" style="width: 1px; height: 1px; visibility: hidden;"><option value="regular">Группа</option><option value="masterclass">Мастер-класс</option><option value="intensive">Интенсив</option><option value="demo">Обучение сотрудников</option><option value="individual">Индивидуальная</option></select></td><td data-col-seq="status"><div class="kv-plugin-loading loading-groupsearch-status">&nbsp;</div><input type="hidden" name="GroupSearch[status]" value=""><select id="groupsearch-status" class="form-control" name="GroupSearch[status][]" multiple size="4" data-s2-options="s2options_3feda8aa" data-krajee-select2="select2_0a4032cd" style="width: 1px; height: 1px; visibility: hidden;"><option value="active" selected>Активная</option><option value="not_started">Не стартовала</option><option value="recruiting">Идет набор</option><option value="suspended">Приостановлена</option><option value="inactive">Окончена</option><option value="destroyed">Развалилась</option></select></td><td data-col-seq="is_online"><div class="kv-plugin-loading loading-groupsearch-is_online">&nbsp;</div><input type="hidden" name="GroupSearch[is_online]" value=""><select id="groupsearch-is_online" class="form-control" name="GroupSearch[is_online][]" multiple size="4" data-s2-options="s2options_3feda8aa" data-krajee-select2="select2_1b859d7b" style="width: 1px; height: 1px; visibility: hidden;"><option value="online">Онлайн</option><option value="offline">Офлайн</option></select></td></tr></thead><tbody><tr class="group-grid" data-key="98637162"><td class="group-grid" data-col-seq="id">98637162</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98637162" data-pjax="0">Библиотека 7 вс 14.00</a><p class="note">Группа по курсу КГ</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">9&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">19.01.2025&nbsp;14:00</td><td class="group-grid" data-col-seq="nextLessonNumber">17</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0ad6e42-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">Знакомство с презентациями</a><div class="note">КГ М5У1</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98623409"><td class="group-grid" data-col-seq="id">98623409</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98623409" data-pjax="0">Библиотека 7 сб 16.00</a><p class="note">Группа по курсу ОЛиП МП</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">5&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">18.01.2025&nbsp;16:00</td><td class="group-grid" data-col-seq="nextLessonNumber">18</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/5923afce-919b-47be-9b17-1154c42a827b" data-pjax="0" data-qa-id="lesson-title">Урок 18. Использование сообщений в игре</a><div class="note">М5У2, курс "Основы логики и программирования", 2021-2022</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98623404"><td class="group-grid" data-col-seq="id">98623404</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98623404" data-pjax="0">Библиотека 7 вс 12.00</a><p class="note">Группа по курсу ОЛиП МП</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">6&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">19.01.2025&nbsp;12:00</td><td class="group-grid" data-col-seq="nextLessonNumber">19</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/5923afce-919b-47be-9b17-1154c42a827b" data-pjax="0" data-qa-id="lesson-title">Урок 18. Использование сообщений в игре</a><div class="note">М5У2, курс "Основы логики и программирования", 2021-2022</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98621252"><td class="group-grid" data-col-seq="id">98621252</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98621252" data-pjax="0">Библиотека 7 вс 18.00</a><p class="note">Группа по курсу Пст</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">9&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">19.01.2025&nbsp;18:00</td><td class="group-grid" data-col-seq="nextLessonNumber">18</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0b6d7fe-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">М5 У1. ООП. Объекты и методы</a><div class="note">Python Start 2021/2022 М5 У1</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98619913"><td class="group-grid" data-col-seq="id">98619913</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98619913" data-pjax="0">Библиотека № 7 сб 10.00</a><p class="note">Группа по курсу КГ</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">9&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">18.01.2025&nbsp;10:00</td><td class="group-grid" data-col-seq="nextLessonNumber">18</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0ad6f81-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">Структурируем презентацию</a><div class="note">КГ М5У2</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98619873"><td class="group-grid" data-col-seq="id">98619873</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98619873" data-pjax="0">Библиотека № 7 сб 14.00</a><p class="note">Группа по курсу ГД</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">8&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">18.01.2025&nbsp;14:00</td><td class="group-grid" data-col-seq="nextLessonNumber">18</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0b6c569-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">Добавление персонажей в игру</a><div class="note">ГД 21/22 М4У2</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98619867"><td class="group-grid" data-col-seq="id">98619867</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98619867" data-pjax="0">Библиотека № 7 сб 12.00</a><p class="note">Группа по курсу ВП</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">9&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">18.01.2025&nbsp;12:00</td><td class="group-grid" data-col-seq="nextLessonNumber">18</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0b67067-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">М4У2. Цикл с условием</a><div class="note">Визуальное программирование, М4У2</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="98589447"><td class="group-grid" data-col-seq="id">98589447</td><td class="group-grid" data-col-seq="title"><a href="/group/view/98589447" data-pjax="0">Библиотека № 7 вс 10.00</a><p class="note">Группа по курсу ВП</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">8&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">19.01.2025&nbsp;10:00</td><td class="group-grid" data-col-seq="nextLessonNumber">19</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0b67067-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">М4У2. Цикл с условием</a><div class="note">Визуальное программирование, М4У2</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="985504"><td class="group-grid" data-col-seq="id">985504</td><td class="group-grid" data-col-seq="title"><a href="/group/view/985504" data-pjax="0">Библиотека 7 сб 18.00</a><p class="note">Группа по курсу Пст 2</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">6&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">18.01.2025&nbsp;18:00</td><td class="group-grid" data-col-seq="nextLessonNumber">19</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0b6e38b-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">М4 У2. Приложение Easy Editor. Ч.1</a><div class="note">Python Start - 2 2021/2022 М4 У2</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr><tr class="group-grid" data-key="978298"><td class="group-grid" data-col-seq="id">978298</td><td class="group-grid" data-col-seq="title"><a href="/group/view/978298" data-pjax="0">Библиотека 7 вс 16.00</a><p class="note">Группа по курсу Пст 2</p></td><td class="group-grid" data-col-seq="venue">Библиотека №7</td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="activeStudentCount"><span title="Учеников в группе / с ноутбуками">7&nbsp;(0)</span></td><td class="group-grid" style="text-align:center; width:1%" data-col-seq="transferredStudentCount">0</td><td class="group-grid" data-col-seq="nextLessonTime">19.01.2025&nbsp;16:00</td><td class="group-grid" data-col-seq="nextLessonNumber">19</td><td class="group-grid" data-col-seq="nextLessonTitle"><a href="/lesson/view/a0c263fa-cde5-11eb-a724-6cb31107bf10" data-pjax="0" data-qa-id="lesson-title">М4 У4. Приложение Easy Editor. Ч. 3</a><div class="note">Python Start - 2 2021/2022 М4 У4</div></td><td class="group-grid" data-col-seq="teacher">Данил Павлов</td><td class="group-grid" data-col-seq="curator">Максим Козлов</td><td class="group-grid" data-col-seq="type"><span class="label label-default">Группа</span></td><td class="group-grid" data-col-seq="status"><span class="label label-success">Активная</span></td><td class="group-grid" data-col-seq="is_online"><div class="label label-info">Офлайн</div></td></tr></tbody></table></div><div class="row"><div class="col-xs-12 text-center"><div class="pagination__wrapper"><ul class="pagination"><li class="active"><a href="/group?GroupSearch%5Bstatus%5D%5B0%5D=active&amp;presetType=all&amp;_pjax=%23group-grid-pjax&amp;page=1">1</a></li></ul></div></div></div></div></div>        </div></div></div></div></div></div></div></div><!-- common modal placeholder --><div class="modal fade" data-backdrop="static" id="common-modal" tabindex="-1" role="dialog" aria-hidden="true"><div class="modal-dialog"><div class="modal-content"><div class="modal-body"><button type="button" class="close pull-right" data-dismiss="modal" aria-label="Закрыть"><span aria-hidden="true">&times;</span></button><ul class="modal-errors"></ul><div class="ajax-content"><h4>Loading..</h4></div></div></div></div></div><script src="/assets/3df76f92/i18nextXHRBackend.min.js?v=1737106652"></script><script src="/assets/687a9c4d/i18nextLocalStorageCache.min.js?v=1737106652"></script></body></html>`

var Kids = []clients.AllGroupsUser{
	{Title: "Группа по курсу КГ", GroupId: "98637162", TimeLesson: "19.01.2025 14:00", RegularTime: "Библиотека 7 вс 14.00"},
	{Title: "Группа по курсу ОЛиП МП", GroupId: "98623409", TimeLesson: "18.01.2025 16:00", RegularTime: "Библиотека 7 сб 16.00"},
	{Title: "Группа по курсу ОЛиП МП", GroupId: "98623404", TimeLesson: "19.01.2025 12:00", RegularTime: "Библиотека 7 вс 12.00"},
	{Title: "Группа по курсу Пст", GroupId: "98621252", TimeLesson: "19.01.2025 18:00", RegularTime: "Библиотека 7 вс 18.00"},
	{Title: "Группа по курсу КГ", GroupId: "98619913", TimeLesson: "18.01.2025 10:00", RegularTime: "Библиотека № 7 сб 10.00"},
	{Title: "Группа по курсу ГД", GroupId: "98619873", TimeLesson: "18.01.2025 14:00", RegularTime: "Библиотека № 7 сб 14.00"},
	{Title: "Группа по курсу ВП", GroupId: "98619867", TimeLesson: "18.01.2025 12:00", RegularTime: "Библиотека № 7 сб 12.00"},
	{Title: "Группа по курсу ВП", GroupId: "98589447", TimeLesson: "19.01.2025 10:00", RegularTime: "Библиотека № 7 вс 10.00"},
	{Title: "Группа по курсу Пст 2", GroupId: "985504", TimeLesson: "18.01.2025 18:00", RegularTime: "Библиотека 7 сб 18.00"},
	{Title: "Группа по курсу Пст 2", GroupId: "978298", TimeLesson: "19.01.2025 16:00", RegularTime: "Библиотека 7 вс 16.00"},
}
