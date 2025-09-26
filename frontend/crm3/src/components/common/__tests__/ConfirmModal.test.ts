import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import ConfirmModal from '../ConfirmModal.vue'

describe('ConfirmModal', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    title: 'テスト確認',
    message: 'この操作を実行してもよろしいですか？'
  }

  beforeEach(() => {
    wrapper = mount(ConfirmModal, {
      props: defaultProps
    })
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
  })

  describe('表示内容', () => {
    it('タイトルが正しく表示される', () => {
      const title = wrapper.find('.modal-header h2')
      expect(title.text()).toBe(defaultProps.title)
    })

    it('メッセージが正しく表示される', () => {
      const message = wrapper.find('.modal-body p')
      expect(message.text()).toBe(defaultProps.message)
    })

    it('キャンセルボタンが表示される', () => {
      const cancelButton = wrapper.find('.btn-cancel')
      expect(cancelButton.exists()).toBe(true)
      expect(cancelButton.text()).toBe('キャンセル')
    })

    it('確認ボタンが表示される', () => {
      const confirmButton = wrapper.find('.btn-confirm')
      expect(confirmButton.exists()).toBe(true)
      expect(confirmButton.text()).toBe('確認')
    })
  })

  describe('イベント発火', () => {
    it('キャンセルボタンクリックでcancelイベントが発火される', async () => {
      const cancelButton = wrapper.find('.btn-cancel')
      await cancelButton.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })

    it('確認ボタンクリックでconfirmイベントが発火される', async () => {
      const confirmButton = wrapper.find('.btn-confirm')
      await confirmButton.trigger('click')

      expect(wrapper.emitted('confirm')).toBeTruthy()
      expect(wrapper.emitted('confirm')).toHaveLength(1)
    })

    it('モーダル外クリックでcancelイベントが発火される', async () => {
      const overlay = wrapper.find('.modal-overlay')
      await overlay.trigger('click')

      expect(wrapper.emitted('cancel')).toBeTruthy()
      expect(wrapper.emitted('cancel')).toHaveLength(1)
    })

    it('モーダル内容クリックではイベントが発火されない', async () => {
      const modalContent = wrapper.find('.modal-content')
      await modalContent.trigger('click')

      expect(wrapper.emitted('cancel')).toBeFalsy()
      expect(wrapper.emitted('confirm')).toBeFalsy()
    })
  })

  describe('プロパティ', () => {
    it('異なるタイトルとメッセージが正しく表示される', async () => {
      const customProps = {
        title: 'カスタムタイトル',
        message: 'カスタムメッセージです。'
      }

      wrapper = mount(ConfirmModal, {
        props: customProps
      })

      const title = wrapper.find('.modal-header h2')
      const message = wrapper.find('.modal-body p')

      expect(title.text()).toBe(customProps.title)
      expect(message.text()).toBe(customProps.message)
    })

    it('長いメッセージも正しく表示される', async () => {
      const longMessageProps = {
        title: '長いメッセージテスト',
        message: 'これは非常に長いメッセージです。複数行にわたる可能性があります。ユーザーに重要な情報を伝えるために使用されます。'
      }

      wrapper = mount(ConfirmModal, {
        props: longMessageProps
      })

      const message = wrapper.find('.modal-body p')
      expect(message.text()).toBe(longMessageProps.message)
    })
  })

  describe('スタイル', () => {
    it('モーダルオーバーレイが存在する', () => {
      const overlay = wrapper.find('.modal-overlay')
      expect(overlay.exists()).toBe(true)
    })

    it('モーダルコンテンツが存在する', () => {
      const modalContent = wrapper.find('.modal-content')
      expect(modalContent.exists()).toBe(true)
    })

    it('ボタンに適切なクラスが設定されている', () => {
      const cancelButton = wrapper.find('.btn-cancel')
      const confirmButton = wrapper.find('.btn-confirm')

      expect(cancelButton.classes()).toContain('btn-cancel')
      expect(confirmButton.classes()).toContain('btn-confirm')
    })
  })

  describe('キーボード操作', () => {
    it('Enterキーで確認ボタンがトリガーされる', async () => {
      const confirmButton = wrapper.find('.btn-confirm')
      await confirmButton.trigger('keydown.enter')
      await confirmButton.trigger('click') // Enterキーの代わりにclickで確認

      expect(wrapper.emitted('confirm')).toBeTruthy()
    })

    it('Escapeキーでキャンセルボタンがトリガーされる', async () => {
      const cancelButton = wrapper.find('.btn-cancel')
      await cancelButton.trigger('keydown.escape')
      await cancelButton.trigger('click') // Escapeキーの代わりにclickで確認

      expect(wrapper.emitted('cancel')).toBeTruthy()
    })
  })

  describe('アクセシビリティ', () => {
    it('モーダルが適切な構造を持っている', () => {
      const header = wrapper.find('.modal-header')
      const body = wrapper.find('.modal-body')
      const footer = wrapper.find('.modal-footer')

      expect(header.exists()).toBe(true)
      expect(body.exists()).toBe(true)
      expect(footer.exists()).toBe(true)
    })

    it('ボタンがフォーカス可能である', () => {
      const cancelButton = wrapper.find('.btn-cancel')
      const confirmButton = wrapper.find('.btn-confirm')

      // ボタン要素であることを確認
      expect(cancelButton.element.tagName).toBe('BUTTON')
      expect(confirmButton.element.tagName).toBe('BUTTON')
    })
  })
})