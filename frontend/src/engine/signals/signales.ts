
export class Signal<T = void> {
	private listeners: ((data: T) => void)[] = [];

	public Connect(listener: (data: T) => void) {
		this.listeners.push(listener);
	}

	public Disconnect(listener: (data: T) => void) {
		this.listeners = this.listeners.filter(l => l !== listener);
	}

	public Emit(data: T) {
		this.listeners.forEach(l => l(data));
	}
}
