// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {conf} from '../models';
import {backend} from '../models';
import {sarama} from '../models';

export function AppConf(arg1:string):Promise<conf.Config>;

export function ConnsOpened():Promise<{[key: string]: backend.Conn}>;

export function DBConf(arg1:string):Promise<{[key: string]: conf.KafkaConfig}>;

export function DescribeConsumerGroups(arg1:string,arg2:Array<string>):Promise<{[key: string]: Array<number>}>;

export function DescribeTopics(arg1:string,arg2:Array<string>):Promise<Array<sarama.TopicMetadata>>;

export function GetConsumerGroupMessageLag(arg1:string,arg2:string,arg3:string):Promise<Array<backend.MessageLag>>;

export function GetMessageByOffset(arg1:string,arg2:string,arg3:number):Promise<{[key: number]: Array<backend.MessageEvent>}>;

export function Greet(arg1:string):Promise<string>;

export function ListConsumerGroups(arg1:string):Promise<{[key: string]: string}>;

export function ListTopics(arg1:string):Promise<{[key: string]: sarama.TopicDetail}>;

export function OpenConn(arg1:string):Promise<void>;

export function SetAPPConf(arg1:conf.Config):Promise<void>;

export function SetDBConf(arg1:{[key: string]: conf.KafkaConfig}):Promise<void>;

export function TailMessage(arg1:string,arg2:string):Promise<{[key: number]: Array<sarama.ConsumerMessage>}>;
