package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ovh/cds/sdk"
)

// PublishProjectEvent publish application event
func PublishProjectEvent(ctx context.Context, payload interface{}, key string, u sdk.Identifiable) {
	bts, _ := json.Marshal(payload)
	event := sdk.Event{
		Timestamp:  time.Now(),
		Hostname:   hostname,
		CDSName:    cdsname,
		EventType:  fmt.Sprintf("%T", payload),
		Payload:    bts,
		ProjectKey: key,
	}
	if u != nil {
		event.UserMail = u.GetEmail()
		event.Username = u.GetUsername()
	}
	publishEvent(ctx, event)
}

// PublishAddProject publishes an event for the creation of the given project
func PublishAddProject(ctx context.Context, p *sdk.Project, u sdk.Identifiable) {
	e := sdk.EventProjectAdd{
		Variables:   p.Variables,
		Permissions: p.ProjectGroups,
		Keys:        p.Keys,
		Metadata:    p.Metadata,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishUpdateProject publishes an event for the modification of the project
func PublishUpdateProject(ctx context.Context, p *sdk.Project, oldProject *sdk.Project, u sdk.Identifiable) {
	e := sdk.EventProjectUpdate{
		NewName:     p.Name,
		NewMetadata: p.Metadata,
		OldMetadata: oldProject.Metadata,
		OldName:     oldProject.Name,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteProject publishess an event for the deletion of the given project
func PublishDeleteProject(ctx context.Context, p *sdk.Project, u sdk.Identifiable) {
	e := sdk.EventProjectDelete{}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishAddProjectVariable publishes an event for the creation of the given variable
func PublishAddProjectVariable(ctx context.Context, p *sdk.Project, v sdk.ProjectVariable, u sdk.Identifiable) {
	if sdk.NeedPlaceholder(v.Type) {
		v.Value = sdk.PasswordPlaceholder
	}
	e := sdk.EventProjectVariableAdd{
		Variable: v,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishUpdateProjectVariable publishes an event for the modification of a variable
func PublishUpdateProjectVariable(ctx context.Context, p *sdk.Project, newVar sdk.ProjectVariable, oldVar sdk.ProjectVariable, u sdk.Identifiable) {
	if sdk.NeedPlaceholder(newVar.Type) {
		newVar.Value = sdk.PasswordPlaceholder
	}
	if sdk.NeedPlaceholder(oldVar.Type) {
		oldVar.Value = sdk.PasswordPlaceholder
	}
	e := sdk.EventProjectVariableUpdate{
		NewVariable: newVar,
		OldVariable: oldVar,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteProjectVariable publishes an event on project variable deletion
func PublishDeleteProjectVariable(ctx context.Context, p *sdk.Project, v sdk.ProjectVariable, u sdk.Identifiable) {
	if sdk.NeedPlaceholder(v.Type) {
		v.Value = sdk.PasswordPlaceholder
	}
	e := sdk.EventProjectVariableDelete{
		Variable: v,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishAddProjectPermission publishes an event on adding a group permission on the project
func PublishAddProjectPermission(ctx context.Context, p *sdk.Project, gp sdk.GroupPermission, u sdk.Identifiable) {
	e := sdk.EventProjectPermissionAdd{
		Permission: gp,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishUpdateProjectPermission publishes an event on updating a group permission on the project
func PublishUpdateProjectPermission(ctx context.Context, p *sdk.Project, gp sdk.GroupPermission, oldGP sdk.GroupPermission, u sdk.Identifiable) {
	e := sdk.EventProjectPermissionUpdate{
		NewPermission: gp,
		OldPermission: oldGP,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteProjectPermission publishes an event on deleting a group permission on the project
func PublishDeleteProjectPermission(ctx context.Context, p *sdk.Project, gp sdk.GroupPermission, u sdk.Identifiable) {
	e := sdk.EventProjectPermissionDelete{
		Permission: gp,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishAddProjectKey publishes an event on adding a project key
func PublishAddProjectKey(ctx context.Context, p *sdk.Project, k sdk.ProjectKey, u sdk.Identifiable) {
	k.Private = sdk.PasswordPlaceholder
	e := sdk.EventProjectKeyAdd{
		Key: k,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

func PublishDisableProjectKey(ctx context.Context, p *sdk.Project, k sdk.ProjectKey, u sdk.Identifiable) {
	k.Private = sdk.PasswordPlaceholder
	e := sdk.EventProjectKeyDisable{
		Key: k,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

func PublishEnableProjectKey(ctx context.Context, p *sdk.Project, k sdk.ProjectKey, u sdk.Identifiable) {
	k.Private = sdk.PasswordPlaceholder
	e := sdk.EventProjectKeyEnable{
		Key: k,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteProjectKey publishes an event on deleting a project key
func PublishDeleteProjectKey(ctx context.Context, p *sdk.Project, k sdk.ProjectKey, u sdk.Identifiable) {
	k.Private = sdk.PasswordPlaceholder
	e := sdk.EventProjectKeyDelete{
		Key: k,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishAddVCSServer publishes an event on adding a project server
func PublishAddVCSServer(ctx context.Context, p *sdk.Project, vcsServerName string, u sdk.Identifiable) {
	e := sdk.EventProjectVCSServerAdd{
		VCSServerName: vcsServerName,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteVCSServer publishes an event on deleting a project server
func PublishDeleteVCSServer(ctx context.Context, p *sdk.Project, vcsServerName string, u sdk.Identifiable) {
	e := sdk.EventProjectVCSServerDelete{
		VCSServerName: vcsServerName,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishAddProjectIntegration publishes an event on adding a integration
func PublishAddProjectIntegration(ctx context.Context, p *sdk.Project, pf sdk.ProjectIntegration, u sdk.Identifiable) {
	pf.Blur()
	e := sdk.EventProjectIntegrationAdd{
		Integration: pf,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishUpdateProjectIntegration publishes an event on updating a integration
func PublishUpdateProjectIntegration(ctx context.Context, p *sdk.Project, pf sdk.ProjectIntegration, pfOld sdk.ProjectIntegration, u sdk.Identifiable) {
	pf.Blur()
	pfOld.Blur()
	e := sdk.EventProjectIntegrationUpdate{
		NewsIntegration: pf,
		OldIntegration:  pfOld,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

// PublishDeleteProjectIntegration publishes an event on deleting integration
func PublishDeleteProjectIntegration(ctx context.Context, p *sdk.Project, pf sdk.ProjectIntegration, u sdk.Identifiable) {
	e := sdk.EventProjectIntegrationDelete{
		Integration: pf,
	}
	PublishProjectEvent(ctx, e, p.Key, u)
}

func PublishRemoveProjectRepository(ctx context.Context, pKey string, vcs sdk.VCSProject, r sdk.ProjectRepository, u sdk.Identifiable) {
	e := sdk.EventProjectRepositoryDelete{
		Repository: r,
		VCS:        vcs,
	}
	PublishProjectEvent(ctx, e, pKey, u)
}

func PublishAddProjectRepository(ctx context.Context, pKey string, vcs sdk.VCSProject, r sdk.ProjectRepository, u sdk.Identifiable) {
	e := sdk.EventProjectRepositoryAdd{
		Repository: r,
		VCS:        vcs,
	}
	PublishProjectEvent(ctx, e, pKey, u)
}

func PublishProjectRepositoryAnalyze(ctx context.Context, pKey string, vcsID string, repoID string, analysisID string, analysisStatus string) {
	e := sdk.EventProjectRepositoryAnalyze{
		RepositoryID: repoID,
		VCSID:        vcsID,
		AnalysisID:   analysisID,
		Status:       analysisStatus,
	}
	PublishProjectEvent(ctx, e, pKey, nil)
}
