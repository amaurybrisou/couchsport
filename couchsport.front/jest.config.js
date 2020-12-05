module.exports = {
  testURL: 'http://localhost',
  moduleFileExtensions: ['js', 'json', 'vue'],

  transform: {
    '^.+\\.js$': '<rootDir>/node_modules/babel-jest',
    '^.*\\.vue$': '<rootDir>/node_modules/vue-jest'
  },

  //   collectCoverage: true,
  //   collectCoverageFrom: ['**/src/*.{js,vue}', '!**/node_modules/**']
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
    '^actions/(.*)$': '<rootDir>/src/store/actions/$1',
    '^plugins/(.*)$': '<rootDir>/src/plugins/$1',
    '^mixins/(.*)$': '<rootDir>/src/mixins/$1',
    '^repos/(.*)$': '<rootDir>/src/repositories/$1'
  },

  preset: '@vue/cli-plugin-unit-jest'
}
