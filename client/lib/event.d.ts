import { Globular } from './services';
/**
 * That local and distant event hub.
 */
export declare class EventHub {
    readonly service: any;
    readonly subscribers: any;
    readonly subscriptions: any;
    readonly globular: Globular;
    /**
     * @param {*} service If undefined only local event will be allow.
     */
    constructor(service: any, globular: Globular);
    /**
     * @param {*} name The name of the event to subcribe to.
     * @param {*} onsubscribe That function return the uuid of the subscriber.
     * @param {*} onevent That function is call when the event is use.
     */
    subscribe(name: string, onsubscribe: (uuid: string) => any, onevent: () => any): void;
    /**
     *
     * @param {*} name
     * @param {*} uuid
     */
    unSubscribe(name: string, uuid: string): void;
    /**
     * Publish an event on the bus, or locally in case of local event.
     * @param {*} name The  name of the event to publish
     * @param {*} data The data associated with the event
     * @param {*} local If the event is not local the data must be seraliaze before sent.
     */
    publish(name: string, data: any, local: boolean): void;
    /** Dispatch the event localy */
    dispatch(name: string, data: any): void;
}
