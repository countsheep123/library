export default {
	data() {
		return {
		};
	},
	methods: {
		login() {
			window.location.href="/oauth2/sign_in";
		},
		loadImage(url) {
			const img=require(`@/assets${url}`);
			return img;
		},
	},
};
