// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "strconv"

func Dashboard(repositories []string, theme string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"space-y-6\"><header><h1 class=\"text-2xl font-bold\">Dashboard</h1></header><div class=\"bg-white dark:bg-neutral-800 rounded-lg shadow p-6\"><div class=\"grid grid-cols-1 md:grid-cols-3 gap-4\"><div class=\"bg-neutral-100 dark:bg-neutral-700 rounded-lg p-4\"><h3 class=\"text-lg font-medium mb-2\">Total Repositories</h3><p class=\"text-3xl font-bold\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(len(repositories)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/pages/dashboard.templ`, Line: 14, Col: 68}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</p></div><div class=\"bg-neutral-100 dark:bg-neutral-700 rounded-lg p-4\"><h3 class=\"text-lg font-medium mb-2\">Registry Status</h3><p class=\"text-green-500 font-medium\">Connected</p></div><div class=\"bg-neutral-100 dark:bg-neutral-700 rounded-lg p-4\"><h3 class=\"text-lg font-medium mb-2\">Quick Actions</h3><div><a href=\"/images\" class=\"text-blue-600 dark:text-blue-400 hover:underline\">Browse Images</a></div></div></div></div><div><h2 class=\"text-xl font-semibold mb-4\">Recent Repositories</h2><div hx-get=\"/htmx/images\" hx-trigger=\"load\" hx-swap=\"outerHTML\"><div class=\"animate-pulse grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i := 0; i < 6; i++ {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<div class=\"bg-neutral-100 dark:bg-neutral-800 rounded-lg p-4 h-16\"></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
