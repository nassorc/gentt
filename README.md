# Entity Component System (ECS)
---
ECS library inspired by donburi's api and entt's sparse set implementation.

This library is a toy project and by no means an alternative to existing libraries such as Donburi. 
This is a second iteration of a ecs library I built a while back, and the goal is to build a new library
that is much simpler and better than the first one, which I felt was unnecessarily complicated.

## Sparse sets
---
Components are stored in a sparse set. 

## Paging
---

## Registering components
Registering the components of your game is reactive. Instead of having to manually register the components,
it happens automatically when the user creates an entity. However, creating the component type is still
done manually, since this will be the handle used to access the store.


Credits to Donburi, EnTT, and the article by skypjack https://skypjack.github.io/2019-02-14-ecs-baf-part-1/.
