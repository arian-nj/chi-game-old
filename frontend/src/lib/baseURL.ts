
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
  const mode = import.meta.env.MODE
  if (mode == "production") {
    return "https://api.chigame.site"
  }
  return "http://localhost:8383";
}
