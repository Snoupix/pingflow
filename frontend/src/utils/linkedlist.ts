export class Node<T> {
	public next: Node<T> | null = null;

	constructor(public data: T) {}
}

export class LinkedList<T> {
	constructor(public head: Node<T> | null = null) {}

	public prepend(node: Node<T>) {
		let head = node;
		while (head.next != null) {
			head = head.next;
		}

		head.next = this.head;
		this.head = head;
	}

	public append(node: Node<T>) {
		if (this.head == null) {
			this.head = node;
			return;
		}

		let head = this.head;
		while (head.next != null) {
			head = head.next;
		}

		head.next = node;
	}

	public to_array(): Array<T> {
		const array: Array<T> = [];
		if (!this.head) {
			return array;
		}

		const array_push = (node: Node<T>): Array<T> => {
			array.push(node.data);
			return node.next == null ? array : array_push(node.next);
		};
		return array_push(this.head);
	}
}

// const linkedlist = new LinkedList();
// linkedlist.prepend(new Node(5));
// linkedlist.append(new Node(1));
// linkedlist.prepend(new Node(10));
// console.log(linkedlist);
// const linkedlist2 = new LinkedList();
// linkedlist2.append(new Node(5));
// linkedlist2.prepend(new Node(1));
// linkedlist2.append(new Node(10));
// console.log(linkedlist2);
