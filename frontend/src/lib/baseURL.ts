
// const apiUrl = import.meta.env.VITE_API_URL

// let url;
// switch (apiUrl) {
// 	case 'production':
// 		url = 'https://stackoverflow.com';
// 		break;
// 	case 'development':
// 	default:
// 		url = 'https://google.com';
// }
// console.log("base url is ", apiUrl)
//
export const GetApiUrl = () => {
  return import.meta.env.VITE_API_URL || "http://localhost:8383";
}
