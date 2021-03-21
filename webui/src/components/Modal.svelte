<script>
  import { createEventDispatcher } from 'svelte'
  import { fly } from 'svelte/transition'

  import ClickOutside from 'svelte-click-outside'

  const dispatch = createEventDispatcher()

  function close() {
    dispatch('close')
  }
</script>

<div
  transition:fly
  class="absolute top-0 left-0 flex items-center justify-center w-full h-full"
  style="background-color: rgba(0,0,0,.4);">
  <ClickOutside on:clickoutside={close}>
    <div
      class="h-auto p-4 mx-2 text-left bg-white rounded shadow-md md:p-6 lg:p-8 md:mx-0">
      <div class="flex flex-col space-y-8">
        <slot name="header" />
        <slot name="main" />
      </div>
    </div>
  </ClickOutside>
</div>
