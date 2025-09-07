export class Message {
	text: string
	userID: bigint
	constructor(text: string, userID: bigint) {
		this.text = text
		this.userID = userID
	}
}

