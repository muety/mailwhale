<script>
  import { onMount } from 'svelte'

  import Layout from '../layouts/Main.svelte'
  import Navigation from '../components/Navigation.svelte'
  import Modal from '../components/Modal.svelte'
  import { createTemplate, deleteTemplate, getDefaultTemplateContent, getTemplates, updateTemplate, } from '../api/template'
  import { errors, successes } from '../stores/alerts'
  import { extractVars } from '../utils/template'

  let templates = []
  let loading = false

  const emptyTemplate = {
    id: null,
    name: '',
    content: '',
  }

  let newTemplate
  let currentTemplate

  let previewVars
  let previewContent

  let newTemplateModal
  let previewModal

  reset()

  async function _createTemplate() {
    try {
      const result = await createTemplate(newTemplate)
      templates = [...templates, result]
      currentTemplate = result
      successes.spawn('Template created successfully')
    } finally {
      newTemplateModal = false
      reset()
    }
  }

  async function _deleteTemplate(id) {
    try {
      await deleteTemplate(id)
      templates = templates.filter((c) => c.id !== id)
      successes.spawn('Template deleted successfully')
    } finally {
      reset()
    }
  }

  async function _updateTemplate() {
    try {
      const result = await updateTemplate(currentTemplate)
      templates[templates.findIndex((t) => t.id === result.id)] = result
      successes.spawn('Template updated successfully')
    } catch (e) {
    }
  }

  async function _loadDefault() {
    try {
      currentTemplate.content = await getDefaultTemplateContent()
      successes.spawn('Template contents replaced by default. Save manually.')
    } catch (e) {
    }
  }

  function reset() {
    previewVars = '{}'
    newTemplate = JSON.parse(JSON.stringify(emptyTemplate))
    if (templates.length) currentTemplate = templates[0]
    else currentTemplate = JSON.parse(JSON.stringify(emptyTemplate))
  }

  function renderHtml(content, vars) {
    let parsedVars
    try {
      parsedVars = JSON.parse(vars)
    } catch (e) {
      errors.spawn('Invalid JSON input')
    }
    let rendered = content
    for (const key in parsedVars) {
      if (parsedVars[key] === null) continue
      rendered = rendered.replace(
        new RegExp(`{{\s*${key.replace(/\./, '\\.')}\s*}}`),
        parsedVars[key]
      )
    }
    return rendered
  }

  function togglePreviewDialog() {
    previewModal = !previewModal
    if (previewModal) {
      previewVars = JSON.stringify(
        extractVars(currentTemplate.content),
        null,
        4
      )
    }
  }

  onMount(async () => {
    loading = true
    try {
      templates = await getTemplates()
    } finally {
      reset()
      loading = false
    }
  })
</script>

<style scoped>
</style>

<Layout>
  <div slot="content" class="flex">
    <div class="w-1/4">
      <Navigation/>
    </div>
    <div class="flex flex-col px-12 w-full w-3/4">
      <div class="flex justify-between mb-8">
        <h1 class="text-2xl font-semibold">Manage Templates</h1>
        <button
          class="flex items-center px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark"
          on:click|stopPropagation={(e) => (newTemplateModal = true)}><span
          class="material-icons">add</span>
          Create
        </button>
      </div>

      <p class="mb-8 ">
        <span class="material-icons" style="font-size: inherit;">info</span>
        Here you can manage mail templates, which are stored in the system and
        can later be referenced when sending a mail. Using templates, you do not
        have to pass the entire mail content from your client application.<br/><br/>Templates
        will usually consist of styled HTML, but can also be plain text. They
        can contain placeholder variables, which are then filled from a JSON
        object when requesting to send a mail using the respective template.<br><br>
        When sending a mail, reference your template using <span class="font-mono text-sm">template_id</span> and <span class="font-mono text-sm">template_vars</span>. A basic, responsive <a href="https://github.com/leemunroe/responsive-html-email-template" target="_blank" rel="noreferrer" class="text-primary">default template</a> is included and can be loaded as a starter.
      </p>

      {#if templates.length}
        <form
          class="flex flex-col space-y-4"
          on:submit|preventDefault={_updateTemplate}>
          <div class="flex items-center w-full mr-2 space-x-4">
            <select
              class="border-2 border-primary rounded-md p-2 flex-grow cursor-pointer"
              bind:value={currentTemplate}>
              {#each templates as template}
                <option value={template}>{template.name} ({template.id})</option>
              {/each}
            </select>
            <button
              type="button"
              class="flex items-center px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
              on:click|stopPropagation={confirm('Are you sure you want to delete this template?') && _deleteTemplate(currentTemplate.id)}><span
              class="material-icons">delete_outline</span></button>
          </div>

          <div>
            <span class="font-semibold">Template ID: </span>
            <span class="font-mono">{currentTemplate.id}</span>
          </div>

          <div class="flex w-full h-full">
            <textarea
              class="font-mono flex-grow w-full text-sm border-2 border-primary rounded-md p-2 flex-grow"
              placeholder="Template content goes here. This can be HTML or plain text."
              style="min-height: 400px"
              bind:value={currentTemplate.content}/>
          </div>

          <div class="flex justify-between w-full">
            <div class="flex space-x-2">
              <span
                type="button"
                class="text-primary underline cursor-pointer"
                on:click|stopPropagation={togglePreviewDialog}>Preview HTML</span>
              <span>|</span>
              <span
                type="button"
                class="text-primary underline cursor-pointer"
                on:click={_loadDefault}>Load Default</span>
            </div>
            <button
              type="submit"
              class="flex items-center px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark"><span
              class="material-icons">save</span>&nbsp; Save
            </button>
          </div>
        </form>
      {:else if loading}
        <div
          class="flex items-center justify-center w-full py-12 text-gray-500">
          <i>Loading ...</i>
        </div>
      {:else}
        <div
          class="w-full py-12 text-gray-500 flex justify-center items-center">
          <i>No templates available. Create your first one.</i>
        </div>
      {/if}
    </div>

    {#if newTemplateModal}
      <Modal on:close={(e) => (newTemplateModal = false) || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">
          Create new template
        </h1>
        <div slot="main" style="min-width: 400px;">
          <form
            class="w-full flex flex-col space-y-4"
            on:submit|preventDefault={_createTemplate}>
            <div class="flex flex-col w-full space-y-1">
              <label for="name-input" class="font-semibold">Name</label>
              <input
                type="text"
                class="border-2 border-primary rounded-md p-2"
                name="name-input"
                placeholder="New template's name"
                required
                bind:value={newTemplate.name}/>
            </div>

            <div class="flex justify-between py-2">
              <div/>
              <button
                type="submit"
                class="py-2 px-4 text-white bg-primary rounded-md hover:bg-primary-dark">Create
              </button>
            </div>
          </form>
        </div>
      </Modal>
    {/if}

    {#if previewModal}
      <Modal on:close={(e) => togglePreviewDialog() || reset()}>
        <h1 class="text-2xl font-semibold" slot="header">Template Preview</h1>
        <div slot="main" style="min-width: 75vw;">
          <div class="flex flex-col space-y-6">
            <div class="flex flex-col space-y-2">
              <label
                for="preview-vars-input"
                class="font-semibold">Variables</label>
              <textarea
                name="preview-vars-input"
                class="font-mono flex-grow w-full text-sm border-2 border-primary rounded-md p-2 flex-grow"
                placeholder="JSON object of variables to be used in the template"
                style="min-height: 250px"
                bind:value={previewVars}/>
              <div class="flex justify-end w-full">
                <button
                  type="button"
                  class="flex items-center px-4 py-2 bg-primary text-white rounded hover:bg-primary-dark"
                  on:click|stopPropagation={(e) => (previewContent = renderHtml(currentTemplate.content, previewVars))}><span
                  class="material-icons">preview</span>
                  &nbsp; Render Preview
                </button>
              </div>
            </div>
            <div class="flex flex-col space-y-2">
              <label for="preview" class="font-semibold">HTML Preview</label>
              <div name="preview" style="min-height: 400px">
                {@html previewContent || currentTemplate.content}
              </div>
            </div>
          </div>
        </div>
      </Modal>
    {/if}
  </div>
</Layout>
