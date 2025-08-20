import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { globalIgnores } from 'eslint/config'

export default tseslint.config([
	globalIgnores(['dist']),
	{
		files: ['**/*.{ts,tsx}'],
		extends: [
			js.configs.recommended,
			tseslint.configs.recommended,
			reactHooks.configs['recommended-latest'],
			reactRefresh.configs.vite,
		],
		languageOptions: {
			ecmaVersion: 2020,
			globals: globals.browser,
		},
		// rules: {
		// 	'import/no-restricted-paths': [
		// 		'error',
		// 		{
		// 			zones: [
		// 				{
		// 					target: './src/features',
		// 					from: './src/app',
		// 				},
		// 				{
		// 					target: [
		// 						'./src/components',
		// 						'./src/lib',
		// 						'./src/types',
		// 						'./src/utils',
		// 					],
		// 					from: ['./src/app', './src/features'],
		// 				},
		// 			],
		// 		}
		// 	]
		// }
	},
])
