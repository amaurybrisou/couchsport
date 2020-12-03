export default {
  data() {
    return {
      help: null
    }
  },
  created: function () {
    const introJS = require('intro.js')
    let i = introJS()
    i.setOptions({
      nextLabel: this.$t('help.next'),
      doneLabel: this.$t('help.done'),
      skipLabel: this.$t('help.skip'),
      prevLabel: this.$t('help.prev'),
      keyboardNavigation: true,
      showBullets: false,
      showStepNumbers: false,
      overlayOpacity: 0.5,
      tooltipClass: 'help',
      highlightClass: 'highlight-help'
    })

    this.help = i
  },
  template: ''
}
