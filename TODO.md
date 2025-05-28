# TODO (Roadmap)

# Road to V1

## V0.1.0 Clean Start
- [X] clean up repo for new purpose.
- [ ] move rag logic to vectorstore for cleaner code reuse between internetsearch
and RAG. This could also be expanded to other mediums down the road.
- [ ] implement simple semantic router. This may have to wait for the previous step.
- [ ] create rough draft for function calling system. This has to couple with the new team system.
- [ ] outline team system.
    - configuration system.
    - debate style of pipeline execpt models can answer with less rigdity.
    - models answer based on routing to save resources.
    - look into multiple go.mods to manage dependencies, or package versions.
- [ ] generate tests based on used resources

## V0.2.0 Rebuild
- [ ] conversation history
- [ ] websockets for streaming data
- [ ] outline distributed setup and plan.
- [ ] look into dynamic promting, models write prompts for models.
- [ ] work out stt and tts system.
- [ ] new cli system. tui? cli?
- [ ] gather hardware information correctly.
    - may have to use a new library or write in house version.
- [ ] change pipelines to new dynamic system.
- [ ] remove docker as container backend logic.
- [ ] move config and folders to cache folders and configs.
    - cleans up configs and models

## V0.3.0 Step Back
- [ ] restucture old code for new systems.
- [ ] improve huggingface download functionality for remote downloading.
- [ ] look into optimizing current code base.
- [ ] look into security for models and app itself.

## V0.4.0

## V0.5.0

## V0.6.0

## V0.7.0
