// originally grifted from, thanks dude
// https://gist.github.com/Phryxia/8f4426db8135fef62744092e0234b399

// export interface SkipListNode<T> {
//     value: T;
//     links: {
//         next?: SkipListNode<T>;
//         prev?: SkipListNode<T>;
//     }[];
// }

// /**
//  * SkipList is a sorted list with randomized algorithm.
//  * Its expected time complexity is O(lg n) for each search.
//  * It eventually resembles binary search tree.
//  * It provides iterator based insertion/deletion.
//  */
// export class SkipList<T, K extends Partial<T> = T> {
//     head: SkipListNode<T> = {
//         value: null as T,
//         links: [],
//     };

//     tail: SkipListNode<T> = {
//         value: null as T,
//         links: [],
//     };

//     private count = 0;

//     /**
//      * @param compare return -1 if a < b, return 0 if a = b, return 1 if a > b
//      * @param maxTopLevel Typically it's lg(n) where n is number of total expected total elements
//      */
//     constructor(
//         private compare: (a: K, b: K) => number,
//         public maxTopLevel: number = Infinity,
//     ) {
//         this.head.links[0] = { next: this.tail };
//         this.tail.links[0] = { prev: this.head };
//     }

//     /**
//      * @param key comparable value with respect to elements
//      * @returns An iterator pointing maximum element equal or less than key, and
//      *          overshoot nodes having minimum element greater than key for each level.
//      * @note overshoots is sorted in decreasing order
//      */
//     findLowerBound(key: K): {
//         iterator: SkipListIterator<T, K>;
//         overshoots: SkipListNode<T>[];
//     } {
//         let node = this.head;
//         let level = this.topLevel();
//         let next = node.links[level].next;
//         const overshoots: SkipListNode<T>[] = [];

//         while (next) {
//             if (
//                 next !== this.tail &&
//                 this.compare(key, next.value as unknown as K) >= 0
//             ) {
//                 node = next;
//                 next = node.links[level].next;
//                 continue;
//             }

//             overshoots.push(next);

//             if (level > 0) {
//                 level -= 1;
//                 next = node.links[level].next;
//             } else {
//                 break;
//             }
//         }
//         return {
//             iterator: new SkipListIterator(this, node),
//             overshoots,
//         };
//     }

//     /**
//      * Add new value to this list
//      * @param value adding target
//      * @returns iterator pointing to the newly added element
//      */
//     add(value: T): SkipListIterator<T, K> {
//         // match new top level
//         const topLevel = this.rollTopLevel();
//         while (topLevel > this.topLevel()) {
//             this.addLevel();
//         }

//         const { iterator, overshoots } = this.findLowerBound(
//             value as unknown as K,
//         );
//         const newNode = { value, links: [] };

//         for (let level = 0; level <= topLevel; ++level) {
//             const upperBound = overshoots[overshoots.length - 1 - level];
//             this.insertBefore(newNode, upperBound, level);
//         }
//         this.count += 1;

//         // since iterator points the lower bound, proceed once
//         iterator.next();
//         return iterator;
//     }

//     /**
//      * Remove value from this list, pointed by given iterator
//      */
//     removeByIterator(itr: SkipListIterator<T, K>): void {
//         const target = itr.node;

//         if (target === this.head || target === this.tail) {
//             throw new Error('Removing head or tail is prohibited');
//         }

//         let emptyLevels = 0;
//         const topLevel = target.links.length - 1;
//         for (let level = 0; level <= topLevel; ++level) {
//             const prev = target.links[level].prev!;
//             const next = target.links[level].next!;

//             prev.links[level].next = next;
//             next.links[level].prev = prev;

//             if (prev === this.head && next === this.tail) {
//                 emptyLevels += 1;
//             }
//         }

//         // preserve only one empty level at top
//         if (emptyLevels >= 2) {
//             this.head.links.length = this.head.links.length - emptyLevels + 1;
//             this.tail.links.length = this.tail.links.length - emptyLevels + 1;
//         }
//     }

//     size(): number {
//         return this.count;
//     }

//     topLevel(): number {
//         return this.head.links.length - 1;
//     }

//     private rollTopLevel(): number {
//         let level = 0;
//         while (level < this.maxTopLevel && Math.random() < 0.5) {
//             level += 1;
//         }
//         return level;
//     }

//     private addLevel() {
//         this.head.links.push({ next: this.tail });
//         this.tail.links.push({ prev: this.head });
//     }

//     private insertBefore(
//         newNode: SkipListNode<T>,
//         base: SkipListNode<T>,
//         level: number,
//     ): void {
//         newNode.links[level] = {
//             prev: base.links[level].prev,
//             next: base,
//         };
//         base.links[level].prev!.links[level].next = newNode;
//         base.links[level].prev = newNode;
//     }
// }

// class SkipListIterator<T, K extends Partial<T>> implements IterableIterator<T> {
//     constructor(public list: SkipList<T, K>, public node: SkipListNode<T>) {}

//     next(): IteratorResult<T> {
//         if (this.node === this.list.tail) {
//             return { done: true, value: undefined };
//         }
//         const value = this.node.value;
//         this.node = this.node.links[0].next!;
//         return { value };
//     }

//     [Symbol.iterator]() {
//         return this;
//     }
// }