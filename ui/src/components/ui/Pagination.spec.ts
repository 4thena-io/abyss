import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import Pagination from './Pagination.vue';

describe('Pagination', () => {
  it('renders nothing when there is only one page', () => {
    const wrapper = mount(Pagination, {
      props: { modelValue: 1, total: 5, pageSize: 10 },
    });
    expect(wrapper.find('div').exists()).toBe(false);
  });

  it('shows the correct item range', () => {
    const wrapper = mount(Pagination, {
      props: { modelValue: 2, total: 25, pageSize: 10 },
    });
    expect(wrapper.text()).toContain('Showing 11 to 20 of 25');
  });

  it('caps the end item at the total on the last page', () => {
    const wrapper = mount(Pagination, {
      props: { modelValue: 3, total: 25, pageSize: 10 },
    });
    expect(wrapper.text()).toContain('Showing 21 to 25 of 25');
  });

  it('disables the previous button on the first page and next on the last', () => {
    const first = mount(Pagination, { props: { modelValue: 1, total: 25, pageSize: 10 } });
    const firstButtons = first.findAll('button');
    const firstPrev = firstButtons[0];
    const firstNext = firstButtons[firstButtons.length - 1];
    expect(firstPrev).toBeDefined();
    expect(firstNext).toBeDefined();
    expect(firstPrev!.attributes('disabled')).toBeDefined();
    expect(firstNext!.attributes('disabled')).toBeUndefined();

    const last = mount(Pagination, { props: { modelValue: 3, total: 25, pageSize: 10 } });
    const lastButtons = last.findAll('button');
    const lastPrev = lastButtons[0];
    const lastNext = lastButtons[lastButtons.length - 1];
    expect(lastPrev).toBeDefined();
    expect(lastNext).toBeDefined();
    expect(lastPrev!.attributes('disabled')).toBeUndefined();
    expect(lastNext!.attributes('disabled')).toBeDefined();
  });

  it('emits update:modelValue with the next page when the next button is clicked', async () => {
    const wrapper = mount(Pagination, {
      props: { modelValue: 1, total: 25, pageSize: 10 },
    });

    const buttons = wrapper.findAll('button');
    const nextButton = buttons[buttons.length - 1];
    expect(nextButton).toBeDefined();
    await nextButton!.trigger('click');

    expect(wrapper.emitted('update:modelValue')).toEqual([[2]]);
  });

  it('emits update:modelValue with the clicked page number', async () => {
    const wrapper = mount(Pagination, {
      props: { modelValue: 1, total: 25, pageSize: 10 },
    });

    const pageButtons = wrapper.findAll('button').filter((b) => b.text() === '3');
    expect(pageButtons[0]).toBeDefined();
    await pageButtons[0]!.trigger('click');

    expect(wrapper.emitted('update:modelValue')).toEqual([[3]]);
  });
});
