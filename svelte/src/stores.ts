import { writable, Writable } from 'svelte/store';

type instance = {
  Name: string;
  ImgFolders: string[];
  TimeToQueue: string;
  PostInterval: string;
  PostDelayAtStartup: string;
  Platforms: {facebook:string, twitter: string}; 
  ItemsInQueue: number;
  NextPostTime: string;
  Status?: string;
}

export const instances: Writable<instance[]> = writable([])
export const activeInstance: Writable<number> = writable(-1) 
