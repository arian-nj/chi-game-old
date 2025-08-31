
const apiUrl = import.meta.env.VITE_API_URL

export const GetBaseUrl = () => {
	// let url;
	// switch (apiUrl) {
	// 	case 'production':
	// 		url = 'https://stackoverflow.com';
	// 		break;
	// 	case 'development':
	// 	default:
	// 		url = 'https://google.com';
	// }

	return apiUrl;
}
