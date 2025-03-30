import IAnimal from './Animal'

type Action = 'Create' | 'Update'

export interface IAnimalModal {
  action: Action,
  title: string,
  data?: IAnimal,
}
