import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import PageTableSection from '../PageTableSection.vue'
import PageToolbar from '../PageToolbar.vue'

describe('phase2 page patterns', () => {
  it('renders toolbar slots with wrap class', () => {
    const wrapper = mount(PageToolbar, {
      slots: {
        left: '<button class="left-action">left</button>',
        right: '<button class="right-action">right</button>',
      },
    })

    expect(wrapper.find('.gi-page-toolbar').exists()).toBe(true)
    expect(wrapper.find('.left-action').exists()).toBe(true)
    expect(wrapper.find('.right-action').exists()).toBe(true)
  })

  it('renders table section shell without inline style', () => {
    const wrapper = mount(PageTableSection, {
      slots: {
        default: '<div class="table-content">table</div>',
      },
    })

    expect(wrapper.find('.gi-table-section').exists()).toBe(true)
    expect(wrapper.html()).not.toContain(' style="')
  })
})
