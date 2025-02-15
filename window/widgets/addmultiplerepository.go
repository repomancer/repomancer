package widgets

//
//type multirepo struct {
//	widget.BaseWidget
//	repositories *ShortcutHandlingEntry
//	output       *widget.Label
//	submitButton *widget.Button
//	cancelButton *widget.Button
//}
//
//func newMultirepo(window fyne.Window) *multirepo {
//	m := &multirepo{
//		repositories: NewShortcutHandlingEntry(window, false),
//		output:       widget.NewLabel(""),
//		submitButton: widget.NewButton("Add All", nil),
//		cancelButton: widget.NewButton("Cancel", nil),
//	}
//	m.ExtendBaseWidget(m)
//	m.repositories.MultiLine = true
//	m.repositories.SetPlaceHolder("github.com/organization/repository")
//	m.repositories.Wrapping = fyne.TextWrapOff
//	m.repositories.OnChanged = func(s string) {
//		m.validate()
//	}
//
//	m.output.Wrapping = fyne.TextWrapWord
//	m.submitButton.Importance = widget.HighImportance
//	m.submitButton.Disable()
//
//	return m
//}
//
//func (m *multirepo) getRepositories() []string {
//	var output []string
//	for _, line := range strings.Split(m.repositories.Text, "\n") {
//		if strings.TrimSpace(line) == "" { // Skip blank lines
//			continue
//		}
//		output = append(output, line)
//	}
//	return output
//}
//
//func (m *multirepo) validate() bool {
//	var errors []string
//	for _, line := range m.getRepositories() {
//		normalized, err := internal.NormalizeGitUrl(line)
//		if err != nil {
//			errors = append(errors, err.Error())
//		}
//		for i := 0; i < project.RepositoryCount(); i++ {
//			existing := project.GetRepository(i)
//			if strings.ToLower(existing.Title()) == strings.ToLower(normalized) {
//				errors = append(errors, fmt.Sprintf("Repository %s already exists in project: %s", line, existing.Title()))
//			}
//		}
//	}
//
//	if len(errors) == 0 {
//		m.submitButton.Enable()
//		m.output.SetText("")
//		return true
//	} else {
//		m.output.SetText(strings.Join(errors, "\n"))
//		m.submitButton.Disable()
//		return false
//	}
//}
//
//func (item *multirepo) CreateRenderer() fyne.WidgetRenderer {
//	form := widget.NewForm(widget.NewFormItem("Repositories", item.repositories))
//	c := container.NewBorder(form, container.NewGridWithColumns(2, item.submitButton, item.cancelButton), nil, nil, item.output)
//
//	return widget.NewSimpleRenderer(c)
//}
//
//func ShowAddMultipleRepositoryWindow(project *internal.Project) {
//	window := NewAppWindow("Add Multiple Repositories", false)
//
//	w := newMultirepo(window)
//
//	w.repositories.OnSubmitted = func(s string) {
//		log.Printf("I was submitted: %s", s)
//	}
//
//	w.submitButton.OnTapped = func() {
//		repos := w.getRepositories()
//		var failed []string
//		var errors []string
//
//		for _, repo := range repos {
//			normalized, err := internal.NormalizeGitUrl(repo)
//			if err != nil {
//				failed = append(failed, repo)
//				e := fmt.Sprintf("%s: %s", repo, err)
//				errors = append(errors, e)
//				w.output.SetText(e)
//				continue
//			}
//
//			w.output.SetText(fmt.Sprintf("Cloning %s and switching to branch %s...", normalized, project.Name))
//			info, err := internal.GetRepositoryInfo(repo)
//			if err != nil {
//				failed = append(failed, repo)
//				e := fmt.Sprintf("%s: %s", repo, err)
//				errors = append(errors, e)
//				w.output.SetText(e)
//				continue
//			}
//
//			err = project.AddRepositoryFromUrl(info.URL)
//			if err != nil {
//				failed = append(failed, repo)
//				e := fmt.Sprintf("%s: %s", repo, err)
//				errors = append(errors, e)
//				w.output.SetText(e)
//				continue
//			}
//			RefreshRepoList()
//		}
//
//		if len(failed) > 0 {
//			w.repositories.SetText(strings.Join(failed, "\n"))
//			w.output.SetText(strings.Join(errors, "\n"))
//		} else {
//			RefreshRepoList()
//			window.Close()
//		}
//	}
//
//	w.cancelButton.OnTapped = func() { window.Close() }
//	window.SetContent(w)
//
//	window.Resize(fyne.NewSize(500, 400))
//	window.Canvas().Focus(w.repositories)
//	window.Show()
//}
