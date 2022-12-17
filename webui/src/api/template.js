import { request } from './api'

async function getDefaultTemplateContent() {
  return (await request('/template/default/content', null, { raw: true })).data
}

async function getTemplates() {
  return (await request('/template', null, {})).data
}

async function renderTemplate(id, data) {
  return (await request(`/template/${ id }/rendered`, data, { method: 'POST' })).data
}

async function createTemplate(template) {
  return (await request('/template', template, { method: 'POST' })).data
}

async function updateTemplate(template) {
  return (await request(`/template/${ template.id }`, template, { method: 'PUT' })).data
}

async function deleteTemplate(id) {
  return await request(`/template/${ id }`, null, { method: 'DELETE' })
}

export { getTemplates, renderTemplate, createTemplate, updateTemplate, deleteTemplate, getDefaultTemplateContent }
