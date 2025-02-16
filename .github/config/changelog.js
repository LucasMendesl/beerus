'use strict'
const config = require('conventional-changelog-conventionalcommits');

module.exports = config({
    "types": [
        { type: 'feat', section: '🚀 New features and improvements' },
        { type: 'fix', section: '🐛 Bug fixes' },
        { type: 'chore', hidden: true },
        { type: 'docs', hidden: true },
        { type: 'style', hidden: true },
        { type: 'refactor', section: '👻 Maintenance' },
        { type: 'perf', section: '👻 Maintenance' },
        { type: 'test', section: '🚦 Tests' },
    ]
})
