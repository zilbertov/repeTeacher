(() => {
  const pages = document.querySelectorAll("[data-page]");
  const navLinks = document.querySelectorAll("[data-route]");
  const title = document.getElementById("page-title");
  const toast = document.getElementById("toast");
  let toastTimer;
  let selectedChatId = null;
  let selectedStudentId = null;
  let selectedTutorId = null;
  let selectedDemoStudentId = null;
  let selectedLessonId = null;
  let selectedNotificationId = null;
  let currentRole = localStorage.getItem("repeTeacherRole") || "tutor";
  let studentsState = [];
  let tutorsState = [];
  let lessonsState = [];
  let chatsState = [];
  let notificationsState = [];
  let profileState = null;
  let notificationSettingsState = null;
  let authLoginPromise = null;
  let selectedCalendarDate = "2026-04-02";
  let calendarMonth = new Date(2026, 3, 1);
  let profileSubjects = [];

  const pageTitles = {
    calendar: "Календарь",
    students: "Ученики",
    notifications: "Уведомления",
    messenger: "Мессенджер",
    settings: "Настройки",
    profile: "Профиль",
  };

  const api = {
    users: "http://127.0.0.1:8081/api",
    lessons: "http://127.0.0.1:8082/api",
  };

  const demoLogin = {
    tutor: { role: "tutor", email: "v4bem@ya.ru", password: "demo" },
    student: { role: "student", email: "student.demo@example.com", password: "demo" },
  };

  const calendarTimes = [
    "09:00",
    "10:00",
    "11:00",
    "12:00",
    "13:00",
    "14:00",
    "15:00",
    "16:00",
    "17:00",
    "18:00",
    "19:00",
    "20:00",
    "21:00",
  ];
  const weekDayNames = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];
  const monthNames = [
    "Январь",
    "Февраль",
    "Март",
    "Апрель",
    "Май",
    "Июнь",
    "Июль",
    "Август",
    "Сентябрь",
    "Октябрь",
    "Ноябрь",
    "Декабрь",
  ];

  if (currentRole !== "tutor" && currentRole !== "student") {
    currentRole = "tutor";
  }

  function isStudentRole() {
    return currentRole === "student";
  }

  function applyRoleClass() {
    document.body.classList.toggle("role-student", isStudentRole());
    document.body.classList.toggle("role-tutor", !isStudentRole());
  }

  function roleText(tutorText, studentText) {
    return isStudentRole() ? studentText : tutorText;
  }

  function roleQuery() {
    return isStudentRole() && selectedDemoStudentId ? `?student_id=${selectedDemoStudentId}` : "";
  }

  function showToast(text) {
    if (!toast || !text) return;
    clearTimeout(toastTimer);
    toast.textContent = text;
    toast.classList.add("is-visible");
    toastTimer = setTimeout(() => toast.classList.remove("is-visible"), 2600);
  }

  function setPage(pageName) {
    const nextPage = pageTitles[pageName] ? pageName : "calendar";
    pages.forEach((page) => {
      page.classList.toggle("is-active", page.dataset.page === nextPage);
    });
    navLinks.forEach((link) => {
      link.classList.toggle("is-active", link.dataset.route === nextPage);
    });
    if (title) title.textContent = nextPage === "students" ? roleText("Ученики", "Репетиторы") : pageTitles[nextPage];
    if (location.hash !== `#${nextPage}`) {
      history.replaceState(null, "", `#${nextPage}`);
    }
  }

  function textSearch(containerId, query) {
    const container = document.getElementById(containerId);
    if (!container) return;
    const value = query.trim().toLowerCase();
    container.querySelectorAll("[data-filter-item]").forEach((item) => {
      item.hidden = !item.textContent.toLowerCase().includes(value);
    });
  }

  function initials(name) {
    return name
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0])
      .join("")
      .toUpperCase();
  }

  function escapeHTML(value) {
    return String(value ?? "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;");
  }

  function setElementValue(element, value) {
    if (!element) return;
    if ("value" in element) {
      element.value = value ?? "";
      return;
    }
    element.textContent = value ?? "";
  }

  function formatTime(value) {
    if (!value) return "";
    return value.slice(0, 5);
  }

  function parseDate(value) {
    const [year, month, day] = value.split("-").map(Number);
    return new Date(year, month - 1, day);
  }

  function formatDate(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }

  function addDays(date, days) {
    const copy = new Date(date);
    copy.setDate(copy.getDate() + days);
    return copy;
  }

  function startOfWeek(date) {
    const copy = new Date(date);
    const dayIndex = (copy.getDay() + 6) % 7;
    copy.setDate(copy.getDate() - dayIndex);
    return copy;
  }

  function getProfileSubjects() {
    const subjects = Array.from(document.querySelectorAll(".subjects-line .subject-pill"))
      .map((item) => item.dataset.profileSubject || item.textContent.replace("×", "").trim())
      .filter(Boolean);
    return subjects.length > 0 ? subjects : ["Математика", "Русский язык"];
  }

  function defaultProfile() {
    const nameElement = document.getElementById("settings-profile-name");
    const emailElement = document.getElementById("settings-profile-email");
    const phoneElement = document.getElementById("settings-profile-phone");
    if (isStudentRole()) {
      return {
        id: selectedDemoStudentId || 0,
        name: (nameElement?.value || nameElement?.textContent || "").trim() || "Тестовый Ученик",
        email: (emailElement?.value || emailElement?.textContent || "").trim() || "student.demo@example.com",
        phone: (phoneElement?.value || phoneElement?.textContent || "").trim() || "89000000000",
        subjects: getProfileSubjects(),
      };
    }
    return {
      id: 1,
      name: (nameElement?.value || nameElement?.textContent || "").trim() || "Вадим Зильбертов",
      email: (emailElement?.value || emailElement?.textContent || "").trim() || "v4bem@ya.ru",
      phone: (phoneElement?.value || phoneElement?.textContent || "").trim() || "89198318673",
      subjects: getProfileSubjects(),
    };
  }

  function closeDialog(form) {
    if (!form) return;
    const dialog = form.closest ? form.closest("dialog") : null;
    if (dialog) dialog.close();
  }

  function updateCalendarSelection(date) {
    selectedCalendarDate = date;
    const parsedDate = parseDate(date);
    if (
      parsedDate.getMonth() !== calendarMonth.getMonth() ||
      parsedDate.getFullYear() !== calendarMonth.getFullYear()
    ) {
      calendarMonth = new Date(parsedDate.getFullYear(), parsedDate.getMonth(), 1);
    }
    const dateInput = document.querySelector("#add-lesson-form [name='lesson_date']");
    if (dateInput) dateInput.value = date;
    renderCalendar();
    renderSelectedDaySchedule();
  }

  function renderProfile(profile) {
    profileState = profile || defaultProfile();
    profileSubjects = Array.isArray(profileState.subjects) ? profileState.subjects : [];

    const settingsName = document.getElementById("settings-profile-name");
    const settingsPhone = document.getElementById("settings-profile-phone");
    const settingsEmail = document.getElementById("settings-profile-email");
    const profileTitle = document.getElementById("profile-title");
    const profilePhone = document.getElementById("profile-phone");
    const profileEmail = document.getElementById("profile-email");
    const profileSubjectsText = document.getElementById("profile-subjects-text");
    const userChipName = document.getElementById("user-chip-name");
    const userChipRole = document.getElementById("user-chip-role");
    const profileRoleLabel = document.getElementById("profile-role-label");
    const roleEyebrow = document.getElementById("role-eyebrow");

    setElementValue(settingsName, profileState.name);
    setElementValue(settingsPhone, profileState.phone);
    setElementValue(settingsEmail, profileState.email);
    if (profileTitle) profileTitle.textContent = profileState.name;
    if (profilePhone) profilePhone.textContent = profileState.phone;
    if (profileEmail) profileEmail.textContent = profileState.email;
    if (profileSubjectsText) profileSubjectsText.textContent = profileSubjects.join(", ") || "Предметы не указаны";
    if (userChipName) userChipName.textContent = profileState.name;
    if (userChipRole) userChipRole.textContent = roleText("Репетитор", "Ученик");
    if (profileRoleLabel) profileRoleLabel.textContent = roleText("Репетитор", "Ученик");
    if (roleEyebrow) roleEyebrow.textContent = roleText("Личный кабинет репетитора", "Личный кабинет ученика");

    const subjectsLine = document.getElementById("profile-subjects");
    if (subjectsLine) {
      subjectsLine.innerHTML = [
        '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M4 4.5A2.5 2.5 0 0 1 6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5Z"/></svg>',
        ...profileSubjects.map(
          (subject) => `<span class="subject-pill" data-profile-subject="${escapeHTML(subject)}">${escapeHTML(subject)} <button type="button" aria-label="Убрать предмет" data-remove-subject>×</button></span>`
        ),
      ].join("");
    }

    renderSubjectOptions();
  }

  function renderSubjectOptions() {
    profileSubjects = profileState?.subjects || getProfileSubjects();
    const lessonSubject = document.querySelector("#add-lesson-form [name='subject']");
    if (lessonSubject) {
      lessonSubject.innerHTML =
        profileSubjects.length > 0
          ? profileSubjects.map((subject) => `<option value="${escapeHTML(subject)}">${escapeHTML(subject)}</option>`).join("")
          : '<option value="">Нет предметов в профиле</option>';
    }

    const studentSubjects = document.getElementById("student-subject-options");
    if (studentSubjects) {
      studentSubjects.innerHTML =
        profileSubjects.length > 0
          ? profileSubjects
              .map(
                (subject, index) => `
            <label>
              <input type="checkbox" name="subjects" value="${escapeHTML(subject)}" ${index === 0 ? "checked" : ""}>
              ${escapeHTML(subject)}
            </label>`
              )
              .join("")
          : '<small class="form-help">Сначала добавьте предмет в профиле.</small>';
    }

    const tutorSubjects = document.getElementById("tutor-subject-options");
    if (tutorSubjects && tutorSubjects.children.length === 0) {
      const defaults = profileSubjects.length > 0 ? profileSubjects : ["Математика", "Русский язык", "Английский язык"];
      tutorSubjects.innerHTML = defaults
        .map(
          (subject, index) => `
            <label>
              <input type="checkbox" name="subjects" value="${escapeHTML(subject)}" ${index === 0 ? "checked" : ""}>
              ${escapeHTML(subject)}
            </label>`
        )
        .join("");
    }
  }

  function renderNotificationSettings(settings) {
    notificationSettingsState =
      settings || notificationSettingsState || {
        push_enabled: true,
        telegram_enabled: true,
        sound_enabled: true,
        lesson_reminders_enabled: false,
      };

    const fields = {
      push_enabled: document.getElementById("notify-push"),
      telegram_enabled: document.getElementById("notify-telegram"),
      sound_enabled: document.getElementById("notify-sound"),
      lesson_reminders_enabled: document.getElementById("notify-lessons"),
    };

    Object.entries(fields).forEach(([key, input]) => {
      if (input) input.checked = Boolean(notificationSettingsState[key]);
    });
  }

  function findLesson(date, time) {
    return lessonsState.find(
      (lesson) =>
        lesson.lesson_date === date &&
        lesson.status !== "cancelled" &&
        formatTime(lesson.start_time).startsWith(time.slice(0, 2))
    );
  }

  function renderLessonCard(lesson) {
    const format = lesson.format === "offline" ? "Очно" : "Онлайн";
    return `
      <button class="lesson-card" type="button" data-open-dialog="lessonDialog" data-lesson-id="${lesson.id}">
        <strong>${escapeHTML(lesson.subject.slice(0, 4))}.</strong>
        <span>${escapeHTML(lesson.exam_type)}</span>
        <span>${format}</span>
        <span>${lesson.has_homework ? "Есть д/з" : ""}</span>
        <small>${lesson.price} р</small>
      </button>`;
  }

  function renderCalendar() {
    const weekGrid = document.getElementById("week-grid");
    const miniCalendar = document.getElementById("mini-calendar");
    const miniTitle = document.getElementById("mini-calendar-title");
    const weekLabel = document.getElementById("week-label");
    const selectedDate = parseDate(selectedCalendarDate);
    const weekStart = startOfWeek(selectedDate);
    const weekDays = Array.from({ length: 7 }, (_, index) => addDays(weekStart, index));

    if (weekLabel) {
      const first = weekDays[0];
      const last = weekDays[6];
      weekLabel.textContent = `Неделя ${first.getDate()} ${monthNames[first.getMonth()].toLowerCase()} - ${last.getDate()} ${monthNames[last.getMonth()].toLowerCase()}`;
    }

    if (weekGrid) {
      const header = [
        '<div class="time-cell week-head"></div>',
        ...weekDays.map((date, index) => {
          const value = formatDate(date);
          return `<button class="day-head ${value === selectedCalendarDate ? "is-selected" : ""}" type="button" data-calendar-date="${value}">${String(date.getDate()).padStart(2, "0")}<br><span>${weekDayNames[index]}</span></button>`;
        }),
      ].join("");

      const rows = calendarTimes
        .map((time) => {
          const slots = weekDays
            .map((date) => {
              const value = formatDate(date);
              const lesson = findLesson(value, time);
              return `<div class="slot">${lesson ? renderLessonCard(lesson) : ""}</div>`;
            })
            .join("");
          return `<div class="time-cell">${time}</div>${slots}`;
        })
        .join("");

      weekGrid.innerHTML = header + rows;
    }

    if (miniTitle) {
      miniTitle.textContent = `${monthNames[calendarMonth.getMonth()]} ${calendarMonth.getFullYear()}`;
    }

    if (miniCalendar) {
      const firstDay = new Date(calendarMonth.getFullYear(), calendarMonth.getMonth(), 1);
      const gridStart = startOfWeek(firstDay);
      const labels = weekDayNames.map((name) => `<span>${name}</span>`).join("");
      const days = Array.from({ length: 42 }, (_, index) => addDays(gridStart, index))
        .map((date) => {
          const value = formatDate(date);
          const muted = date.getMonth() !== calendarMonth.getMonth();
          const selected = value === selectedCalendarDate;
          return `<button type="button" class="${muted ? "muted" : ""} ${selected ? "is-selected" : ""}" data-calendar-date="${value}">${date.getDate()}</button>`;
        })
        .join("");
      miniCalendar.innerHTML = labels + days;
    }
  }

  function renderSelectedDaySchedule() {
    const scheduleRow = document.querySelector(".schedule-row");
    if (!scheduleRow) return;
    const lessons = lessonsState.filter(
      (lesson) => lesson.lesson_date === selectedCalendarDate && lesson.status !== "cancelled"
    );
    if (lessons.length === 0) {
      scheduleRow.disabled = true;
      delete scheduleRow.dataset.lessonId;
      scheduleRow.innerHTML = `<span>Нет занятий на эту дату</span>`;
      return;
    }
    const lesson = lessons[0];
    scheduleRow.disabled = false;
    scheduleRow.dataset.lessonId = String(lesson.id);
    scheduleRow.innerHTML = `
      <span>${formatTime(lesson.start_time)}</span>
      <i aria-hidden="true">•</i>
      <strong>${escapeHTML(lesson.subject)}</strong>
      <small></small>`;
  }

  function updateLessonStudentSelect() {
    const select = document.querySelector("#add-lesson-form [name='student_id']");
    if (!select) return;
    const label = document.getElementById("lesson-person-label");
    if (label) label.textContent = roleText("Ученик", "Репетитор");
    if (isStudentRole()) {
      if (tutorsState.length === 0) {
        select.innerHTML = '<option value="">Сначала добавьте репетитора</option>';
        return;
      }
      select.innerHTML = tutorsState
        .map((tutor) => `<option value="${tutor.id}">${escapeHTML(tutor.name)}</option>`)
        .join("");
      return;
    }
    const students = studentsState.filter((student) => student.status !== "archived");
    if (students.length === 0) {
      select.innerHTML = '<option value="">Сначала добавьте ученика</option>';
      return;
    }
    select.innerHTML = students
      .map((student) => `<option value="${student.id}">${escapeHTML(student.name)}</option>`)
      .join("");
  }

  function selectStudent(row) {
    selectedTutorId = null;
    selectedStudentId = Number(row.dataset.id);
    document.querySelectorAll("[data-select-student]").forEach((item) => {
      item.classList.toggle("is-selected", item === row);
    });
    document.getElementById("student-name").textContent = row.dataset.name;
    document.getElementById("student-email").textContent = row.dataset.email;
    document.getElementById("student-phone").textContent = row.dataset.phone;
    document.getElementById("student-subject").textContent = row.dataset.subject;
    document.getElementById("student-status").textContent = row.dataset.status;
    document.getElementById("student-notes").value = row.dataset.notes || "";
    document.getElementById("student-notes").readOnly = false;
    document.getElementById("student-avatar").textContent = initials(row.dataset.name);
  }

  function selectTutor(row) {
    selectedStudentId = null;
    selectedTutorId = Number(row.dataset.id);
    document.querySelectorAll("[data-select-tutor]").forEach((item) => {
      item.classList.toggle("is-selected", item === row);
    });
    document.getElementById("student-name").textContent = row.dataset.name;
    document.getElementById("student-email").textContent = row.dataset.email;
    document.getElementById("student-phone").textContent = row.dataset.phone;
    document.getElementById("student-subject").textContent = row.dataset.subject;
    document.getElementById("student-status").textContent = "demo";
    document.getElementById("student-notes").value = row.dataset.notes || "";
    document.getElementById("student-notes").readOnly = false;
    document.getElementById("student-avatar").textContent = initials(row.dataset.name);
  }

  function setPeopleModeUI() {
    const studentsTitle = document.getElementById("students-title");
    const peopleEyebrow = document.getElementById("people-eyebrow");
    const peopleSearchLabel = document.getElementById("people-search-label");
    const peopleSearch = document.getElementById("people-search");
    const personRoleLabel = document.getElementById("person-role-label");
    const messengerEyebrow = document.getElementById("messenger-eyebrow");
    const chatInput = document.getElementById("chat-message-input");
    const notes = document.getElementById("student-notes");
    const studentsNavLink = document.querySelector('[data-route="students"]');
    const studentsNav = studentsNavLink?.querySelector("span");
    const saveNotesButton = document.querySelector("[data-save-student-notes]");

    applyRoleClass();

    if (studentsTitle) studentsTitle.textContent = roleText("Ученики", "Репетиторы");
    if (studentsNav) studentsNav.textContent = roleText("Ученики", "Репетиторы");
    if (studentsNavLink) studentsNavLink.title = roleText("Ученики", "Репетиторы");
    if (peopleEyebrow) peopleEyebrow.textContent = roleText("Заявки, активные и архив", "Учебные репетиторы");
    if (peopleSearchLabel) peopleSearchLabel.textContent = roleText("Поиск ученика", "Поиск репетитора");
    if (peopleSearch) peopleSearch.placeholder = roleText("Поиск ученика...", "Поиск репетитора...");
    if (personRoleLabel) personRoleLabel.textContent = roleText("Ученик", "Репетитор");
    if (messengerEyebrow) messengerEyebrow.textContent = roleText("Переписка с учениками", "Переписка с репетиторами");
    if (chatInput) chatInput.placeholder = roleText("Написать сообщение...", "Написать репетитору...");
    if (notes) notes.readOnly = false;
    if (title && location.hash.replace("#", "") === "students") title.textContent = roleText("Ученики", "Репетиторы");

    document.querySelectorAll("[data-archive-student], [data-delete-student]").forEach((button) => {
      button.hidden = isStudentRole();
    });
    document.querySelectorAll('[data-open-dialog="addLessonDialog"]').forEach((button) => {
      button.hidden = isStudentRole();
    });
    if (saveNotesButton) {
      saveNotesButton.hidden = false;
      saveNotesButton.textContent = "Сохранить заметку";
    }
  }

  function setNotificationIcon(kind) {
    const icon = document.getElementById("notification-icon");
    if (!icon) return;
    icon.classList.toggle("danger", kind === "reschedule" || kind === "cancel");
    icon.innerHTML =
      kind === "message"
        ? '<path d="M20 21a8 8 0 0 0-16 0M12 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"/>'
        : '<path d="m21.73 18-8-14a2 2 0 0 0-3.46 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4M12 17h.01"/>';
  }

  function setNotificationActions(kind, row) {
    const actions = document.getElementById("notification-actions");
    if (!actions) return;
    if (kind === "message") {
      actions.innerHTML = [
        '<button class="button button-primary" type="button" data-open-notification-chat>Открыть чат</button>',
        '<button class="button" type="button" data-notification-action="read">Прочитано</button>',
      ].join("");
      return;
    }
    if (kind === "lesson") {
      actions.innerHTML = [
        row?.dataset.lessonId ? '<button class="button button-primary" type="button" data-open-notification-lesson>Открыть занятие</button>' : "",
        '<button class="button" type="button" data-notification-action="read">Прочитано</button>',
      ].join("");
      return;
    }
    if (kind === "cancel") {
      actions.innerHTML = [
        '<button class="button button-primary" type="button" data-notification-action="approve">Принять</button>',
        '<button class="button danger" type="button" data-notification-action="reject">Отказать</button>',
        row?.dataset.lessonId ? '<button class="button" type="button" data-open-notification-lesson>Открыть занятие</button>' : "",
      ].join("");
      return;
    }
    actions.innerHTML = [
      '<button class="button button-primary" type="button" data-notification-action="approve">Согласиться</button>',
      '<button class="button danger" type="button" data-notification-action="reject">Отказать</button>',
      row?.dataset.lessonId ? '<button class="button" type="button" data-open-notification-lesson>Открыть занятие</button>' : "",
    ].join("");
  }

  function selectNotification(row) {
    if (!row) return;
    selectedNotificationId = Number(row.dataset.id);
    document.querySelectorAll("[data-select-notification]").forEach((item) => {
      item.classList.toggle("is-selected", item === row);
    });
    document.getElementById("notification-title").textContent = row.dataset.title;
    document.getElementById("notification-text").textContent = row.dataset.text;
    setNotificationIcon(row.dataset.kind);
    setNotificationActions(row.dataset.kind, row);
  }

  function setSettingsTab(tabName) {
    document.querySelectorAll("[data-tab]").forEach((tab) => {
      tab.classList.toggle("is-active", tab.dataset.tab === tabName);
    });
    document.querySelectorAll("[data-tab-panel]").forEach((panel) => {
      panel.classList.toggle("is-active", panel.dataset.tabPanel === tabName);
    });
  }

  function appendOwnMessage(text) {
    const messages = document.getElementById("messages");
    const now = new Date().toLocaleTimeString("ru-RU", {
      hour: "2-digit",
      minute: "2-digit",
    });
    const item = document.createElement("div");
    item.className = "message is-own";
    item.innerHTML = `<p></p><small>${now}</small>`;
    item.querySelector("p").textContent = text;
    messages.appendChild(item);
    messages.scrollTop = messages.scrollHeight;
  }

  function setChatHead(name) {
    const head = document.querySelector(".chat-head");
    if (!head) return;
    const displayName = name || "Нет выбранного чата";
    head.querySelector(".avatar").textContent = name ? initials(name) : "";
    head.querySelector("strong").textContent = displayName;
  }

  function authStorageKey(role = currentRole) {
    return `repeTeacherAuth:${role}`;
  }

  function readAuthSession(role = currentRole) {
    try {
      return JSON.parse(localStorage.getItem(authStorageKey(role)) || "null");
    } catch {
      return null;
    }
  }

  function writeAuthSession(session) {
    if (!session || !session.role) return;
    localStorage.setItem(authStorageKey(session.role), JSON.stringify(session));
    if (session.role === "student" && session.student_id) {
      selectedDemoStudentId = Number(session.student_id);
    }
  }

  async function loginDemoRole(role = currentRole) {
    const credentials = demoLogin[role] || demoLogin.tutor;
    const response = await fetch(`${api.users}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
    });
    if (!response.ok) {
      throw new Error(`Auth error ${response.status}`);
    }
    const session = await response.json();
    writeAuthSession(session);
    return session;
  }

  async function ensureAuthToken() {
    const saved = readAuthSession();
    const now = Math.floor(Date.now() / 1000);
    if (saved?.token && saved.expires_at && saved.expires_at > now + 60) {
      if (saved.role === "student" && saved.student_id) {
        selectedDemoStudentId = Number(saved.student_id);
      }
      return saved;
    }
    if (!authLoginPromise) {
      authLoginPromise = loginDemoRole(currentRole).finally(() => {
        authLoginPromise = null;
      });
    }
    return authLoginPromise;
  }

  async function apiJSON(url, options = {}) {
    const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
    if (!options.skipAuth) {
      const session = await ensureAuthToken();
      headers.Authorization = `Bearer ${session.token}`;
    }
    const { skipAuth, ...fetchOptions } = options;
    const response = await fetch(url, {
      ...fetchOptions,
      headers,
    });
    if (!response.ok) {
      throw new Error(`API error ${response.status}`);
    }
    if (response.status === 204) return null;
    return response.json();
  }

  function openRoleDialog() {
    const dialog = document.getElementById("roleDialog");
    if (dialog && typeof dialog.showModal === "function") {
      dialog.showModal();
    }
  }

  async function switchRole(role) {
    if (role !== "tutor" && role !== "student") return;
    currentRole = role;
    localStorage.setItem("repeTeacherRole", role);
    authLoginPromise = null;
    await ensureAuthToken();
    selectedChatId = null;
    selectedNotificationId = null;
    selectedLessonId = null;
    setPeopleModeUI();
    closeDialog(document.querySelector("#roleDialog form"));
    setPage("calendar");
    await loadBackendData();
    showToast(roleText("Вы вошли как репетитор", "Вы вошли как тестовый ученик"));
  }

  function renderStudentRow(student, selected) {
    const subjects = (student.subjects || []).join(", ");
    const acceptButton =
      student.status === "request"
        ? `<button class="icon-button accept-button" type="button" data-accept-student="${student.id}">✓</button>`
        : "";
    return `
      <div class="list-row student-row ${selected ? "is-selected" : ""}" tabindex="0" role="button"
        data-filter-item data-select-student data-id="${student.id}"
        data-name="${escapeHTML(student.name)}" data-email="${escapeHTML(student.email)}" data-phone="${escapeHTML(student.phone)}"
        data-subject="${escapeHTML(subjects)}" data-exam="${escapeHTML(student.exam_type)}" data-status="${escapeHTML(student.status)}"
        data-notes="${escapeHTML(student.notes || "")}">
        <span class="avatar"></span>
        <div>
          <strong>${escapeHTML(student.name)}</strong>
          <small>${escapeHTML(subjects)}${student.exam_type ? `, ${escapeHTML(student.exam_type)}` : ""}</small>
        </div>
        ${acceptButton}
      </div>`;
  }

  function renderStudentGroup(titleText, students, selectedID) {
    if (students.length === 0) return "";
    return `<p class="list-title">${titleText}</p>${students.map((student) => renderStudentRow(student, student.id === selectedID)).join("")}`;
  }

  function renderTutorRow(tutor, selected) {
    const subjects = (tutor.subjects || []).join(", ");
    return `
      <div class="list-row student-row ${selected ? "is-selected" : ""}" tabindex="0" role="button"
        data-filter-item data-select-tutor data-id="${tutor.id}"
        data-name="${escapeHTML(tutor.name)}" data-email="${escapeHTML(tutor.email)}" data-phone="${escapeHTML(tutor.phone)}"
        data-subject="${escapeHTML(subjects)}" data-notes="${escapeHTML(tutor.notes || "")}">
        <span class="avatar"></span>
        <div>
          <strong>${escapeHTML(tutor.name)}</strong>
          <small>${escapeHTML(subjects || "Предметы не указаны")}</small>
        </div>
      </div>`;
  }

  function renderTutors(tutors) {
    const list = document.getElementById("students-list");
    if (!list) return;
    tutorsState = tutors || [];
    updateLessonStudentSelect();
    setPeopleModeUI();
    if (!tutors || tutors.length === 0) {
      selectedTutorId = null;
      list.innerHTML = `<p class="list-title">Пока нет репетиторов</p>`;
      document.getElementById("student-name").textContent = "Нет выбранного репетитора";
      document.getElementById("student-email").textContent = "-";
      document.getElementById("student-phone").textContent = "-";
      document.getElementById("student-subject").textContent = "-";
      document.getElementById("student-status").textContent = "-";
      document.getElementById("student-notes").value = "";
      return;
    }
    const selected = tutors.find((tutor) => tutor.id === selectedTutorId) || tutors[0];
    selectedTutorId = selected.id;
    list.innerHTML = `<p class="list-title">Мои репетиторы</p>${tutors.map((tutor) => renderTutorRow(tutor, tutor.id === selected.id)).join("")}`;
    selectTutor(list.querySelector(`[data-id="${selected.id}"]`));
  }

  function renderStudents(students) {
    const list = document.getElementById("students-list");
    if (!list) return;
    studentsState = students || [];
    updateLessonStudentSelect();
    setPeopleModeUI();
    if (!students || students.length === 0) {
      selectedStudentId = null;
      selectedDemoStudentId = null;
      list.innerHTML = `<p class="list-title">Пока нет учеников</p>`;
      document.getElementById("student-name").textContent = "Нет выбранного ученика";
      document.getElementById("student-email").textContent = "-";
      document.getElementById("student-phone").textContent = "-";
      document.getElementById("student-subject").textContent = "-";
      document.getElementById("student-status").textContent = "-";
      document.getElementById("student-notes").value = "";
      return;
    }
    const selected =
      students.find((student) => student.id === selectedStudentId) ||
      students.find((student) => student.status === "active") ||
      students[0];
    list.innerHTML = [
      renderStudentGroup("Заявки", students.filter((student) => student.status === "request"), selected.id),
      renderStudentGroup("Мои ученики", students.filter((student) => student.status === "active"), selected.id),
      renderStudentGroup("Архив", students.filter((student) => student.status === "archived"), selected.id),
    ].join("");
    selectStudent(list.querySelector(`[data-id="${selected.id}"]`));
  }

  function renderLessons(lessons) {
    lessonsState = lessons || [];
    renderCalendar();
    renderSelectedDaySchedule();
  }

  function findLessonByID(id) {
    return lessonsState.find((lesson) => lesson.id === Number(id));
  }

  function findStudentByID(id) {
    return studentsState.find((student) => student.id === Number(id));
  }

  function findTutorByID(id) {
    return tutorsState.find((tutor) => tutor.id === Number(id));
  }

  function ensureDemoStudent() {
    const availableStudents = studentsState.filter((student) => student.status !== "archived");
    if (availableStudents.length === 0) {
      selectedDemoStudentId = null;
      return null;
    }
    if (!selectedDemoStudentId || !availableStudents.some((student) => student.id === selectedDemoStudentId)) {
      const seeded = availableStudents.find((student) => student.email === "student.demo@example.com");
      selectedDemoStudentId = (seeded || availableStudents[0]).id;
    }
    return findStudentByID(selectedDemoStudentId);
  }

  function renderStudentProfile() {
    const student = ensureDemoStudent();
    if (!student) {
      renderProfile({
        id: 0,
        name: "Нет тестового ученика",
        email: "-",
        phone: "-",
        subjects: [],
      });
      return;
    }
    renderProfile(student);
  }

  function renderFileList(lesson, type) {
    const files = (lesson.files || []).filter((file) => file.file_type === type);
    if (files.length === 0) return '<span class="file-pill">Нет файлов</span>';
    return files
      .map((file) => `<a class="file-pill" href="${escapeHTML(file.file_path)}" target="_blank" rel="noreferrer">${escapeHTML(file.file_name)}</a>`)
      .join("");
  }

  function renderLessonDialog(lesson) {
    if (!lesson) return;
    selectedLessonId = lesson.id;

    const studentName = lesson.student_name || findStudentByID(lesson.student_id)?.name || "Ученик";
    const tutorName = lesson.tutor_name || findTutorByID(lesson.tutor_id)?.name || "Репетитор";
    const personName = roleText(studentName, tutorName);
    const statusText = lesson.status === "cancelled" ? "Занятие отменено" : "Запланировано";
    const formatText = lesson.format === "offline" ? "Очно" : "Онлайн";

    setElementValue(document.getElementById("lesson-dialog-student"), personName);
    setElementValue(document.getElementById("lesson-dialog-status"), `${roleText("Ученик", "Репетитор")} · ${statusText}`);
    setElementValue(document.getElementById("lesson-dialog-date"), `Дата - ${lesson.lesson_date}`);
    setElementValue(document.getElementById("lesson-dialog-time"), `Время - ${formatTime(lesson.start_time)}`);
    setElementValue(document.getElementById("lesson-dialog-subject"), `Предмет - ${lesson.subject}`);
    setElementValue(document.getElementById("lesson-dialog-exam"), `Экзамен - ${lesson.exam_type || "-"}`);
    setElementValue(document.getElementById("lesson-dialog-format"), `Формат - ${formatText}`);
    setElementValue(document.getElementById("lesson-dialog-price"), `Цена - ${lesson.price} р`);

    const avatar = document.getElementById("lesson-dialog-avatar");
    if (avatar) avatar.textContent = initials(personName);

    const dateInput = document.getElementById("lesson-reschedule-date");
    const timeInput = document.getElementById("lesson-reschedule-time");
    if (dateInput) dateInput.value = lesson.lesson_date;
    if (timeInput) timeInput.value = formatTime(lesson.start_time);

    const materialList = document.getElementById("lesson-material-files");
    const homeworkList = document.getElementById("lesson-homework-files");
    if (materialList) materialList.innerHTML = renderFileList(lesson, "material");
    if (homeworkList) homeworkList.innerHTML = renderFileList(lesson, "homework");

    const cancelButton = document.querySelector("[data-cancel-lesson]");
    if (cancelButton) cancelButton.hidden = isStudentRole();
  }

  function selectedLesson() {
    return findLessonByID(selectedLessonId);
  }

  async function startChatForStudent(studentID) {
    if (!studentID) {
      showToast("Сначала выберите ученика");
      return null;
    }
    const chat = await apiJSON(`${api.users}/chats`, {
      method: "POST",
      body: JSON.stringify({ student_id: Number(studentID) }),
    });
    selectedChatId = chat.id;
    return chat;
  }

  async function startChatForTutor(tutorID) {
    const student = ensureDemoStudent();
    if (!student) {
      showToast("Сначала нужен тестовый ученик");
      return null;
    }
    if (!tutorID) {
      showToast("Сначала выберите репетитора");
      return null;
    }
    const chat = await apiJSON(`${api.users}/chats`, {
      method: "POST",
      body: JSON.stringify({ student_id: Number(student.id), tutor_id: Number(tutorID) }),
    });
    selectedChatId = chat.id;
    return chat;
  }

  async function startChatForSelectedStudent() {
    try {
      const chat = isStudentRole() ? await startChatForTutor(selectedTutorId) : await startChatForStudent(selectedStudentId);
      if (!chat) return;
      setPage("messenger");
      await loadBackendData();
      await loadMessages(selectedChatId);
      showToast("Чат открыт");
    } catch {
      showToast("Не удалось открыть чат");
    }
  }

  async function startChatForSelectedLesson() {
    const lesson = selectedLesson();
    if (!lesson) {
      showToast("Сначала выберите занятие");
      return;
    }
    try {
      const chat = isStudentRole() ? await startChatForTutor(lesson.tutor_id) : await startChatForStudent(lesson.student_id);
      if (!chat) return;
      closeDialog(document.getElementById("lessonDialog"));
      setPage("messenger");
      await loadBackendData();
      await loadMessages(selectedChatId);
      showToast("Чат открыт");
    } catch {
      showToast("Не удалось открыть чат");
    }
  }

  async function rescheduleSelectedLesson() {
    const lesson = selectedLesson();
    if (!lesson) {
      showToast("Сначала выберите занятие");
      return;
    }
    const lessonDate = document.getElementById("lesson-reschedule-date")?.value;
    const startTime = document.getElementById("lesson-reschedule-time")?.value;
    if (!lessonDate || !startTime) {
      showToast("Укажите дату и время");
      return;
    }
    try {
      const updated = await apiJSON(`${api.lessons}/lessons/${lesson.id}/reschedule`, {
        method: "POST",
        body: JSON.stringify({ lesson_date: lessonDate, start_time: startTime, sender_type: currentRole }),
      });
      selectedLessonId = updated.id;
      await loadBackendData();
      updateCalendarSelection(updated.lesson_date);
      renderLessonDialog(updated);
      showToast("Занятие перенесено");
    } catch {
      showToast("Не удалось перенести занятие");
    }
  }

  async function cancelSelectedLesson() {
    if (isStudentRole()) {
      showToast("Ученик не может отменять занятия");
      return;
    }
    const lesson = selectedLesson();
    if (!lesson) {
      showToast("Сначала выберите занятие");
      return;
    }
    if (!window.confirm("Отменить выбранное занятие?")) return;
    try {
      const updated = await apiJSON(`${api.lessons}/lessons/${lesson.id}/cancel`, {
        method: "POST",
        body: JSON.stringify({ sender_type: currentRole }),
      });
      selectedLessonId = updated.id;
      await loadBackendData();
      renderLessonDialog(updated);
      showToast("Занятие отменено");
    } catch {
      showToast("Не удалось отменить занятие");
    }
  }

  async function addFileToSelectedLesson(type) {
    const lesson = selectedLesson();
    if (!lesson) {
      showToast("Сначала выберите занятие");
      return;
    }
    const input = document.querySelector(`[data-file-name="${type}"]`);
    const fileName = input?.value.trim();
    if (!fileName) {
      showToast("Укажите имя файла");
      return;
    }
    try {
      await apiJSON(`${api.lessons}/lessons/${lesson.id}/files`, {
        method: "POST",
        body: JSON.stringify({
          file_type: type,
          file_name: fileName,
          file_path: `/files/${fileName}`,
        }),
      });
      if (input) input.value = "";
      await loadBackendData();
      renderLessonDialog(findLessonByID(lesson.id));
      showToast("Файл добавлен");
    } catch {
      showToast("Не удалось добавить файл");
    }
  }

  async function applyNotificationAction(action) {
    if (!selectedNotificationId) {
      showToast("Сначала выберите уведомление");
      return;
    }
    try {
      await apiJSON(`${api.users}/notifications/${selectedNotificationId}/${action}`, { method: "POST" });
      await loadBackendData();
      showToast(action === "reject" ? "Отказ сохранен" : "Уведомление обработано");
    } catch {
      showToast("Не удалось обработать уведомление");
    }
  }

  async function readAllNotifications() {
    const unread = notificationsState.filter((item) => !item.is_read);
    if (unread.length === 0) {
      showToast("Непрочитанных уведомлений нет");
      return;
    }
    try {
      await Promise.all(unread.map((item) => apiJSON(`${api.users}/notifications/${item.id}/read`, { method: "POST" })));
      await loadBackendData();
      showToast("Все уведомления прочитаны");
    } catch {
      showToast("Не удалось прочитать все уведомления");
    }
  }

  async function openNotificationChat() {
    const row = document.querySelector("[data-select-notification].is-selected");
    const studentID = Number(row?.dataset.studentId);
    const tutorID = Number(row?.dataset.tutorId);
    try {
      const chat = isStudentRole() ? await startChatForTutor(tutorID) : await startChatForStudent(studentID);
      if (!chat) return;
      selectedChatId = chat.id;
      setPage("messenger");
      await applyNotificationAction("read");
      await loadMessages(selectedChatId);
    } catch {
      showToast("Не удалось открыть чат");
    }
  }

  function openNotificationLesson() {
    const row = document.querySelector("[data-select-notification].is-selected");
    const lessonID = Number(row?.dataset.lessonId);
    const lesson = findLessonByID(lessonID);
    if (!lesson) {
      showToast("Занятие не найдено");
      return;
    }
    selectedCalendarDate = lesson.lesson_date;
    selectedLessonId = lesson.id;
    updateCalendarSelection(lesson.lesson_date);
    setPage("calendar");
    const dialog = document.getElementById("lessonDialog");
    renderLessonDialog(lesson);
    if (dialog && typeof dialog.showModal === "function") dialog.showModal();
  }

  function renderNotifications(notifications) {
    const list = document.getElementById("notifications-list");
    if (!list) return;
    notificationsState = notifications || [];
    if (!notifications || notifications.length === 0) {
      selectedNotificationId = null;
      list.innerHTML = `<p class="list-title">Пока нет уведомлений</p>`;
      document.getElementById("notification-title").textContent = "Нет уведомлений";
      document.getElementById("notification-text").textContent = "";
      document.getElementById("notification-actions").innerHTML = "";
      return;
    }

    list.innerHTML = notifications
      .map((item, index) => {
        const danger = item.type === "reschedule" || item.type === "cancel";
        const unread = item.is_read ? "" : " is-unread";
        const icon = danger
          ? `<svg class="row-icon danger" viewBox="0 0 24 24" aria-hidden="true"><path d="m21.73 18-8-14a2 2 0 0 0-3.46 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4M12 17h.01"/></svg>`
          : `<svg class="row-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M20 21a8 8 0 0 0-16 0M12 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"/></svg>`;
        return `
          <button class="list-row notification-row${unread} ${index === 0 ? "is-selected" : ""}" type="button"
            data-filter-item data-select-notification data-id="${item.id}" data-title="${escapeHTML(item.title)}"
            data-text="${escapeHTML(item.description)}" data-kind="${escapeHTML(item.type)}"
            data-tutor-id="${item.tutor_id || ""}" data-student-id="${item.student_id || ""}" data-lesson-id="${item.lesson_id || ""}">
            ${icon}
            <div><strong>${escapeHTML(item.title)}</strong><small>${escapeHTML(item.description)}</small></div>
            <span class="time">${new Date(item.created_at).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}</span>
          </button>`;
      })
      .join("");
    selectNotification(list.querySelector("[data-select-notification]"));
  }

  function renderChats(chats) {
    const list = document.getElementById("chat-list");
    const messages = document.getElementById("messages");
    if (!list || !messages) return;
    chatsState = chats || [];
    if (!chats || chats.length === 0) {
      selectedChatId = null;
      list.innerHTML = `<p class="list-title">Пока нет чатов</p>`;
      messages.innerHTML = `<div class="message"><p>${roleText('Выберите ученика и нажмите "Написать сообщение".', 'Выберите репетитора и нажмите "Написать сообщение".')}</p></div>`;
      setChatHead("");
      return;
    }

    const selectedChat = chats.find((chat) => chat.id === selectedChatId) || chats[0];
    selectedChatId = selectedChat.id;
    setChatHead(selectedChat.participant_name);
    list.innerHTML = chats
      .map((chat, index) => `
        <button class="list-row chat-row ${chat.id === selectedChatId ? "is-selected" : ""}" type="button" data-filter-item data-chat-id="${chat.id}" data-chat-name="${escapeHTML(chat.participant_name)}">
          <span class="avatar"></span>
          <div><strong>${escapeHTML(chat.participant_name)}</strong><small>${escapeHTML(chat.last_message || "Нет сообщений")}</small></div>
          <span class="time">${chat.last_message_time ? new Date(chat.last_message_time).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" }) : ""}</span>
        </button>`)
      .join("");
    loadMessages(selectedChatId, "messages", currentRole);
  }

  function updateStudentInState(student) {
    studentsState = studentsState.map((item) => (item.id === student.id ? student : item));
    renderStudents(studentsState);
  }

  function updateTutorInState(tutor) {
    tutorsState = tutorsState.map((item) => (item.id === tutor.id ? tutor : item));
    renderTutors(tutorsState);
  }

  async function saveSelectedStudentNotes() {
    if (isStudentRole()) {
      await saveSelectedTutorNotes();
      return;
    }
    if (!selectedStudentId) {
      showToast("Сначала выберите ученика");
      return;
    }
    const notes = document.getElementById("student-notes")?.value || "";
    try {
      const student = await apiJSON(`${api.users}/students/${selectedStudentId}/notes`, {
        method: "POST",
        body: JSON.stringify({ notes }),
      });
      updateStudentInState(student);
      showToast("Заметка сохранена");
    } catch {
      showToast("Не удалось сохранить заметку");
    }
  }

  async function saveSelectedTutorNotes() {
    if (!selectedTutorId) {
      showToast("Сначала выберите репетитора");
      return;
    }
    const notes = document.getElementById("student-notes")?.value || "";
    try {
      const tutor = await apiJSON(`${api.users}/tutors/${selectedTutorId}/notes`, {
        method: "POST",
        body: JSON.stringify({ notes }),
      });
      updateTutorInState(tutor);
      showToast("Заметка сохранена");
    } catch {
      showToast("Не удалось сохранить заметку");
    }
  }

  async function deleteSelectedStudent() {
    if (!selectedStudentId) {
      showToast("Сначала выберите ученика");
      return;
    }
    if (!window.confirm("Удалить выбранного ученика?")) return;
    try {
      await apiJSON(`${api.users}/students/${selectedStudentId}`, { method: "DELETE" });
      studentsState = studentsState.filter((student) => student.id !== selectedStudentId);
      selectedStudentId = null;
      renderStudents(studentsState);
      await loadBackendData();
      showToast("Ученик удален");
    } catch {
      showToast("Не удалось удалить ученика");
    }
  }

  async function archiveSelectedStudent() {
    if (!selectedStudentId) {
      showToast("Сначала выберите ученика");
      return;
    }
    try {
      const student = await apiJSON(`${api.users}/students/${selectedStudentId}/archive`, { method: "POST" });
      updateStudentInState(student);
      await loadBackendData();
      showToast("Ученик перенесен в архив");
    } catch {
      showToast("Не удалось перенести ученика в архив");
    }
  }

  async function saveProfileSubjects(subjects) {
    const currentProfile = profileState || defaultProfile();
    try {
      if (isStudentRole()) {
        const student = ensureDemoStudent();
        if (!student) {
          showToast("Нет тестового ученика");
          return;
        }
        const saved = await apiJSON(`${api.users}/students/${student.id}`, {
          method: "PUT",
          body: JSON.stringify({
            name: currentProfile.name,
            email: currentProfile.email,
            phone: currentProfile.phone,
            subjects,
            exam_type: student.exam_type || "ЕГЭ",
            status: student.status || "active",
            notes: student.notes || "",
          }),
        });
        studentsState = studentsState.map((item) => (item.id === saved.id ? saved : item));
        renderProfile(saved);
        showToast("Предметы ученика сохранены");
        return;
      }
      const profile = await apiJSON(`${api.users}/profile`, {
        method: "PUT",
        body: JSON.stringify({
          name: currentProfile.name,
          email: currentProfile.email,
          phone: currentProfile.phone,
          subjects,
        }),
      });
      renderProfile(profile);
      showToast("Предметы профиля сохранены");
    } catch {
      showToast("Не удалось сохранить предметы");
    }
  }

  async function saveProfileForm() {
    const currentProfile = profileState || defaultProfile();
    const name = document.getElementById("settings-profile-name")?.value.trim();
    const email = document.getElementById("settings-profile-email")?.value.trim();
    const phone = document.getElementById("settings-profile-phone")?.value.trim();
    if (!name || !email) {
      showToast("Укажите имя и email");
      return;
    }

    try {
      if (isStudentRole()) {
        const student = ensureDemoStudent();
        if (!student) {
          showToast("Нет тестового ученика");
          return;
        }
        const saved = await apiJSON(`${api.users}/students/${student.id}`, {
          method: "PUT",
          body: JSON.stringify({
            name,
            email,
            phone,
            subjects: currentProfile.subjects || profileSubjects,
            exam_type: student.exam_type || "ЕГЭ",
            status: student.status || "active",
            notes: student.notes || "",
          }),
        });
        studentsState = studentsState.map((item) => (item.id === saved.id ? saved : item));
        renderProfile(saved);
        showToast("Профиль ученика сохранен");
        return;
      }
      const profile = await apiJSON(`${api.users}/profile`, {
        method: "PUT",
        body: JSON.stringify({
          name,
          email,
          phone,
          subjects: currentProfile.subjects || profileSubjects,
        }),
      });
      renderProfile(profile);
      showToast("Профиль сохранен");
    } catch {
      showToast("Не удалось сохранить профиль");
    }
  }

  async function changePassword() {
    const currentPassword = document.getElementById("current-password")?.value || "";
    const newPassword = document.getElementById("new-password")?.value || "";
    const confirmPassword = document.getElementById("confirm-password")?.value || "";
    if (!newPassword && !confirmPassword) return;
    if (isStudentRole()) {
      showToast("Пароль ученика в учебном режиме не используется");
      return;
    }
    if (newPassword !== confirmPassword) {
      showToast("Пароли не совпадают");
      return;
    }
    try {
      await apiJSON(`${api.users}/profile/password`, {
        method: "POST",
        body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
      });
      ["current-password", "new-password", "confirm-password"].forEach((id) => {
        const input = document.getElementById(id);
        if (input) input.value = "";
      });
      showToast("Пароль изменен");
    } catch {
      showToast("Не удалось изменить пароль");
    }
  }

  async function saveNotificationSettings() {
    const settings = {
      push_enabled: Boolean(document.getElementById("notify-push")?.checked),
      telegram_enabled: Boolean(document.getElementById("notify-telegram")?.checked),
      sound_enabled: Boolean(document.getElementById("notify-sound")?.checked),
      lesson_reminders_enabled: Boolean(document.getElementById("notify-lessons")?.checked),
    };
    try {
      if (isStudentRole()) {
        notificationSettingsState = settings;
        localStorage.setItem("repeTeacherStudentNotifySettings", JSON.stringify(settings));
        renderNotificationSettings(settings);
        showToast("Настройки ученика сохранены локально");
        return;
      }
      const saved = await apiJSON(`${api.users}/settings/notifications`, {
        method: "PUT",
        body: JSON.stringify(settings),
      });
      renderNotificationSettings(saved);
      showToast("Настройки уведомлений сохранены");
    } catch {
      showToast("Не удалось сохранить уведомления");
    }
  }

  function renderMessages(messages, boxId = "messages", ownSender = "tutor") {
    const box = document.getElementById(boxId);
    if (!box) return;
    if (!messages || messages.length === 0) {
      box.innerHTML = `<div class="message"><p>Сообщений пока нет.</p></div>`;
      return;
    }
    box.innerHTML = messages
      .map((message) => `
        <div class="message ${message.sender_type === ownSender ? "is-own" : ""}">
          <p>${escapeHTML(message.text)}</p>
          <small>${new Date(message.created_at).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}</small>
        </div>`)
      .join("");
    box.scrollTop = box.scrollHeight;
  }

  async function loadMessages(chatId, boxId = "messages", ownSender = "tutor") {
    try {
      const messages = await apiJSON(`${api.users}/chats/${chatId}/messages`);
      renderMessages(messages, boxId, ownSender);
    } catch {
      showToast("Не удалось загрузить сообщения");
    }
  }

  async function loadBackendData() {
    setPeopleModeUI();

    const baseResults = await Promise.allSettled([
      apiJSON(`${api.users}/profile`),
      apiJSON(`${api.users}/students`),
      apiJSON(`${api.users}/tutors`),
      apiJSON(`${api.users}/settings/notifications`),
    ]);

    const [profile, students, tutors, settings] = baseResults;
    let loadedCount = 0;

    if (students.status === "fulfilled") {
      loadedCount += 1;
      studentsState = students.value || [];
      ensureDemoStudent();
    }
    if (tutors.status === "fulfilled") {
      loadedCount += 1;
      tutorsState = tutors.value || [];
    }
    if (profile.status === "fulfilled" && !isStudentRole()) {
      loadedCount += 1;
      renderProfile(profile.value);
    }
    if (isStudentRole()) {
      renderStudentProfile();
      renderTutors(tutorsState);
    } else {
      renderStudents(studentsState);
    }
    if (settings.status === "fulfilled" && !isStudentRole()) {
      loadedCount += 1;
      renderNotificationSettings(settings.value);
    } else if (isStudentRole()) {
      const saved = localStorage.getItem("repeTeacherStudentNotifySettings");
      renderNotificationSettings(saved ? JSON.parse(saved) : notificationSettingsState);
    }

    const suffix = roleQuery();
    const roleResults = await Promise.allSettled([
      apiJSON(`${api.lessons}/lessons${suffix}`),
      apiJSON(`${api.users}/notifications${suffix}`),
      apiJSON(`${api.users}/chats${suffix}`),
    ]);

    const [lessons, notifications, chats] = roleResults;
    if (lessons.status === "fulfilled") {
      loadedCount += 1;
      renderLessons(lessons.value);
    }
    if (notifications.status === "fulfilled") {
      loadedCount += 1;
      renderNotifications(notifications.value);
    }
    if (chats.status === "fulfilled") {
      loadedCount += 1;
      renderChats(chats.value);
    }

    updateCalendarSelection(selectedCalendarDate);

    if (loadedCount === 0) {
      showToast("Backend не запущен, показан статический интерфейс");
    } else if (loadedCount < baseResults.length + roleResults.length - (isStudentRole() ? 2 : 0)) {
      showToast("Часть данных не загрузилась");
    }
  }

  document.addEventListener("click", (event) => {
    const dialogButton = event.target.closest("[data-open-dialog]");
    const toastButton = event.target.closest("[data-toast]");
    const gotoButton = event.target.closest("[data-goto-page]");
    const studentRow = event.target.closest("[data-select-student]");
    const tutorRow = event.target.closest("[data-select-tutor]");
    const notificationRow = event.target.closest("[data-select-notification]");
    const tabButton = event.target.closest("[data-tab]");
    const acceptButton = event.target.closest("[data-accept-student]");
    const chatRow = event.target.closest("[data-chat-id]");
    const calendarButton = event.target.closest("[data-calendar-date]");
    const calendarPrev = event.target.closest("[data-calendar-prev]");
    const calendarNext = event.target.closest("[data-calendar-next]");
    const weekPrev = event.target.closest("[data-week-prev]");
    const weekNext = event.target.closest("[data-week-next]");
    const closeDialogButton = event.target.closest("[data-close-dialog]");
    const saveNotesButton = event.target.closest("[data-save-student-notes]");
    const archiveStudentButton = event.target.closest("[data-archive-student]");
    const deleteStudentButton = event.target.closest("[data-delete-student]");
    const startChatButton = event.target.closest("[data-start-chat]");
    const startLessonChatButton = event.target.closest("[data-start-lesson-chat]");
    const rescheduleLessonButton = event.target.closest("[data-reschedule-lesson]");
    const cancelLessonButton = event.target.closest("[data-cancel-lesson]");
    const addLessonFileButton = event.target.closest("[data-add-lesson-file]");
    const notificationActionButton = event.target.closest("[data-notification-action]");
    const readAllNotificationsButton = event.target.closest("[data-read-all-notifications]");
    const openNotificationChatButton = event.target.closest("[data-open-notification-chat]");
    const openNotificationLessonButton = event.target.closest("[data-open-notification-lesson]");
    const saveProfileButton = event.target.closest("[data-save-profile]");
    const saveNotifyButton = event.target.closest("[data-save-notification-settings]");
    const removeSubjectButton = event.target.closest("[data-remove-subject]");
    const roleLogoutButton = event.target.closest("[data-role-logout]");
    const loginRoleButton = event.target.closest("[data-login-role]");

    if (closeDialogButton) {
      closeDialog(closeDialogButton);
    }

    if (calendarButton) {
      updateCalendarSelection(calendarButton.dataset.calendarDate);
    }

    if (calendarPrev || calendarNext) {
      const step = calendarNext ? 1 : -1;
      const nextMonth = new Date(calendarMonth.getFullYear(), calendarMonth.getMonth() + step, 1);
      updateCalendarSelection(formatDate(nextMonth));
    }

    if (weekPrev || weekNext) {
      const step = weekNext ? 7 : -7;
      updateCalendarSelection(formatDate(addDays(parseDate(selectedCalendarDate), step)));
    }

    if (dialogButton) {
      const dialogID =
        dialogButton.dataset.openDialog === "addStudentDialog" && isStudentRole()
          ? "addTutorDialog"
          : dialogButton.dataset.openDialog;
      if (dialogID === "addLessonDialog" && isStudentRole()) {
        showToast("Ученик не может создавать занятия");
        return;
      }
      const dialog = document.getElementById(dialogID);
      if (dialogID === "addLessonDialog") {
        updateLessonStudentSelect();
        updateCalendarSelection(selectedCalendarDate);
      }
      if (dialogID === "lessonDialog") {
        const lessonID = Number(dialogButton.dataset.lessonId);
        const lesson = findLessonByID(lessonID);
        if (lesson) renderLessonDialog(lesson);
      }
      if (dialog && typeof dialog.showModal === "function") dialog.showModal();
    }

    if (loginRoleButton) {
      switchRole(loginRoleButton.dataset.loginRole);
    } else if (roleLogoutButton) {
      openRoleDialog();
    } else if (acceptButton) {
      apiJSON(`${api.users}/students/${acceptButton.dataset.acceptStudent}/accept`, { method: "POST" })
        .then(loadBackendData)
        .then(() => showToast("Заявка принята"))
        .catch(() => showToast("Не удалось принять заявку"));
    } else if (saveNotesButton) {
      saveSelectedStudentNotes();
    } else if (archiveStudentButton) {
      archiveSelectedStudent();
    } else if (deleteStudentButton) {
      deleteSelectedStudent();
    } else if (startChatButton) {
      startChatForSelectedStudent();
    } else if (startLessonChatButton) {
      startChatForSelectedLesson();
    } else if (rescheduleLessonButton) {
      rescheduleSelectedLesson();
    } else if (cancelLessonButton) {
      cancelSelectedLesson();
    } else if (addLessonFileButton) {
      addFileToSelectedLesson(addLessonFileButton.dataset.addLessonFile);
    } else if (notificationActionButton) {
      applyNotificationAction(notificationActionButton.dataset.notificationAction);
    } else if (readAllNotificationsButton) {
      readAllNotifications();
    } else if (openNotificationChatButton) {
      openNotificationChat();
    } else if (openNotificationLessonButton) {
      openNotificationLesson();
    } else if (saveProfileButton) {
      saveProfileForm().then(changePassword);
    } else if (saveNotifyButton) {
      saveNotificationSettings();
    } else if (removeSubjectButton) {
      const pill = removeSubjectButton.closest("[data-profile-subject]");
      const subject = pill?.dataset.profileSubject;
      if (subject) {
        saveProfileSubjects(profileSubjects.filter((item) => item !== subject));
      }
    } else if (toastButton) {
      showToast(toastButton.dataset.toast);
    }
    if (gotoButton) setPage(gotoButton.dataset.gotoPage);
    if (studentRow && !event.target.closest(".accept-button")) selectStudent(studentRow);
    if (tutorRow) selectTutor(tutorRow);
    if (notificationRow) selectNotification(notificationRow);
    if (chatRow) {
      selectedChatId = Number(chatRow.dataset.chatId);
      document.querySelectorAll("[data-chat-id]").forEach((row) => row.classList.toggle("is-selected", row === chatRow));
      setChatHead(chatRow.dataset.chatName);
      loadMessages(selectedChatId, "messages", currentRole);
    }
    if (tabButton) setSettingsTab(tabButton.dataset.tab);
  });

  document.addEventListener("keydown", (event) => {
    const studentRow = event.target.closest("[data-select-student]");
    const tutorRow = event.target.closest("[data-select-tutor]");
    if (studentRow && (event.key === "Enter" || event.key === " ")) {
      event.preventDefault();
      selectStudent(studentRow);
    }
    if (tutorRow && (event.key === "Enter" || event.key === " ")) {
      event.preventDefault();
      selectTutor(tutorRow);
    }
  });

  document.querySelectorAll("[data-filter]").forEach((input) => {
    input.addEventListener("input", () => textSearch(input.dataset.filter, input.value));
  });

  document.getElementById("add-student-form")?.addEventListener("submit", (event) => {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    const subjects = data
      .getAll("subjects")
      .map(String)
      .map((subject) => subject.trim())
      .filter(Boolean);
    if (subjects.length === 0) {
      showToast("Укажите хотя бы один предмет");
      return;
    }
    apiJSON(`${api.users}/students`, {
      method: "POST",
      body: JSON.stringify({
        name: String(data.get("name") || "").trim(),
        email: String(data.get("email") || "").trim(),
        phone: String(data.get("phone") || "").trim(),
        subjects,
        exam_type: String(data.get("exam_type") || "").trim(),
        status: String(data.get("status") || "active"),
        notes: String(data.get("notes") || "").trim(),
      }),
    })
      .then((student) => {
        studentsState = [...studentsState, student];
        renderStudents(studentsState);
        form.reset();
        renderSubjectOptions();
        closeDialog(form);
        setPage("students");
        showToast("Ученик добавлен");
        return loadBackendData();
      })
      .catch(() => showToast("Не удалось добавить ученика"));
  });

  document.getElementById("add-tutor-form")?.addEventListener("submit", (event) => {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    const subjects = data
      .getAll("subjects")
      .map(String)
      .map((subject) => subject.trim())
      .filter(Boolean);
    if (subjects.length === 0) {
      showToast("Укажите хотя бы один предмет");
      return;
    }
    apiJSON(`${api.users}/tutors`, {
      method: "POST",
      body: JSON.stringify({
        name: String(data.get("name") || "").trim(),
        email: String(data.get("email") || "").trim(),
        phone: String(data.get("phone") || "").trim(),
        subjects,
      }),
    })
      .then((tutor) => {
        tutorsState = [...tutorsState, tutor];
        selectedTutorId = tutor.id;
        renderTutors(tutorsState);
        form.reset();
        closeDialog(form);
        setPage("students");
        showToast("Репетитор добавлен");
        return loadBackendData();
      })
      .catch(() => showToast("Не удалось добавить репетитора"));
  });

  document.getElementById("add-lesson-form")?.addEventListener("submit", (event) => {
    event.preventDefault();
    if (isStudentRole()) {
      showToast("Ученик не может создавать занятия");
      return;
    }
    const form = event.currentTarget;
    const data = new FormData(form);
    const selectedPersonID = Number(data.get("student_id"));
    if (!selectedPersonID) {
      showToast(roleText("Сначала выберите ученика", "Сначала выберите репетитора"));
      return;
    }
    const studentID = selectedPersonID;
    const tutorID = 0;
    if (!studentID) {
      showToast("Нет тестового ученика для создания занятия");
      return;
    }
    const lessonDate = String(data.get("lesson_date") || selectedCalendarDate);
    apiJSON(`${api.lessons}/lessons`, {
      method: "POST",
      body: JSON.stringify({
        tutor_id: tutorID,
        student_id: studentID,
        subject: String(data.get("subject") || "").trim(),
        exam_type: String(data.get("exam_type") || "").trim(),
        lesson_date: lessonDate,
        start_time: String(data.get("start_time") || "").trim(),
        duration_minutes: Number(data.get("duration_minutes")),
        format: String(data.get("format") || "online"),
        has_homework: data.has("has_homework"),
        price: Number(data.get("price")),
      }),
    })
      .then((lesson) => {
        lessonsState = [...lessonsState, lesson];
        updateCalendarSelection(lesson.lesson_date);
        closeDialog(form);
        setPage("calendar");
        showToast("Занятие добавлено");
        return loadBackendData();
      })
      .catch(() => showToast("Не удалось добавить занятие"));
  });

  document.querySelector("#add-lesson-form [name='lesson_date']")?.addEventListener("change", (event) => {
    const value = event.currentTarget.value;
    if (value) updateCalendarSelection(value);
  });

  document.getElementById("chat-form")?.addEventListener("submit", (event) => {
    event.preventDefault();
    const input = event.currentTarget.elements.message;
    const text = input.value.trim();
    if (!text) return;
    if (!selectedChatId) {
      showToast("Сначала выберите чат");
      return;
    }
    apiJSON(`${api.users}/chats/${selectedChatId}/messages`, {
      method: "POST",
      body: JSON.stringify({ text, sender_type: currentRole }),
    })
      .then(() => {
        input.value = "";
        return loadBackendData().then(() => loadMessages(selectedChatId, "messages", currentRole));
      })
      .catch(() => {
        appendOwnMessage(text);
        input.value = "";
        showToast("Backend недоступен, сообщение показано только на экране");
      });
  });

  window.addEventListener("hashchange", () => {
    setPage(location.hash.replace("#", ""));
  });

  setPage(location.hash.replace("#", "") || "calendar");
  setPeopleModeUI();
  renderSubjectOptions();
  updateCalendarSelection(selectedCalendarDate);
  loadBackendData();
})();
