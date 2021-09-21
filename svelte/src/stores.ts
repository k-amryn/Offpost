import { writable, Writable } from 'svelte/store';

// instance is read directly from Go
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

// ginstance is converted from instance for the GUI,
// and converted back to instance when user clicks Save
type ginstance = {
  Name: string;
  ImgFolders: string[];
  TimeToQueue: {num: number, unit: string};
  PostInterval: {num: number, unit: string};
  PostDelayAtStartup: string;
  Platforms: {facebook:string, twitter: string}; 
  ItemsInQueue: number;
  NextPostTime: string;
  Status?: string;
}

export const instances: Writable<instance[]> = writable([])
export const ginstances: Writable<ginstance[]> = writable([])
export const activeInstance: Writable<number> = writable(-1) 
