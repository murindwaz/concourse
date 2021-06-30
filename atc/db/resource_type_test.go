package db_test

import (
	"context"
	"time"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/event"
	"github.com/concourse/concourse/tracing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

var _ = Describe("ResourceType", func() {
	var pipeline db.Pipeline

	BeforeEach(func() {
		var (
			created bool
			err     error
		)

		pipeline, created, err = defaultTeam.SavePipeline(
			atc.PipelineRef{Name: "pipeline-with-types"},
			atc.Config{
				ResourceTypes: atc.ResourceTypes{
					{
						Name:     "some-type",
						Type:     "registry-image",
						Source:   atc.Source{"some": "repository"},
						Defaults: atc.Source{"some-default-k1": "some-default-v1"},
					},
					{
						Name:       "some-other-type",
						Type:       "some-type",
						Privileged: true,
						Source:     atc.Source{"some": "other-repository"},
					},
					{
						Name:   "some-type-with-params",
						Type:   "s3",
						Source: atc.Source{"some": "repository"},
						Params: atc.Params{"unpack": "true"},
					},
					{
						Name:       "some-type-with-custom-check",
						Type:       "registry-image",
						Source:     atc.Source{"some": "repository"},
						CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
					},
				},
			},
			0,
			false,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(created).To(BeTrue())
	})

	Describe("(Pipeline).ResourceTypes", func() {
		var resourceTypes db.ResourceTypes

		JustBeforeEach(func() {
			var err error
			resourceTypes, err = pipeline.ResourceTypes()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the resource types", func() {
			Expect(resourceTypes).To(HaveLen(4))

			ids := map[int]struct{}{}

			for _, t := range resourceTypes {
				ids[t.ID()] = struct{}{}

				switch t.Name() {
				case "some-type":
					Expect(t.Name()).To(Equal("some-type"))
					Expect(t.Type()).To(Equal("registry-image"))
					Expect(t.Source()).To(Equal(atc.Source{"some": "repository"}))
					Expect(t.Defaults()).To(Equal(atc.Source{"some-default-k1": "some-default-v1"}))
				case "some-other-type":
					Expect(t.Name()).To(Equal("some-other-type"))
					Expect(t.Type()).To(Equal("some-type"))
					Expect(t.Source()).To(Equal(atc.Source{"some": "other-repository"}))
					Expect(t.Defaults()).To(BeNil())
					Expect(t.Privileged()).To(BeTrue())
				case "some-type-with-params":
					Expect(t.Name()).To(Equal("some-type-with-params"))
					Expect(t.Type()).To(Equal("s3"))
					Expect(t.Source()).To(Equal(atc.Source{"some": "repository"}))
					Expect(t.Defaults()).To(BeNil())
					Expect(t.Params()).To(Equal(atc.Params{"unpack": "true"}))
				case "some-type-with-custom-check":
					Expect(t.Name()).To(Equal("some-type-with-custom-check"))
					Expect(t.Type()).To(Equal("registry-image"))
					Expect(t.Source()).To(Equal(atc.Source{"some": "repository"}))
					Expect(t.Defaults()).To(BeNil())
					Expect(t.CheckEvery().Interval.String()).To(Equal("10ms"))
				}
			}

			Expect(ids).To(HaveLen(4))
		})

		Context("when a resource type becomes inactive", func() {
			BeforeEach(func() {
				var (
					created bool
					err     error
				)

				pipeline, created, err = defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-types"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{
							{
								Name:   "some-type",
								Type:   "registry-image",
								Source: atc.Source{"some": "repository"},
							},
						},
					},
					pipeline.ConfigVersion(),
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeFalse())
			})

			It("does not return inactive resource types", func() {
				Expect(resourceTypes).To(HaveLen(1))
				Expect(resourceTypes[0].Name()).To(Equal("some-type"))
			})
		})

		Context("when the resource types list has resource with same name", func() {
			BeforeEach(func() {
				var (
					created bool
					err     error
				)

				pipeline, created, err = defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-types"},
					atc.Config{
						Resources: atc.ResourceConfigs{
							{
								Name:   "some-name",
								Type:   "some-name",
								Source: atc.Source{},
							},
						},
						ResourceTypes: atc.ResourceTypes{
							{
								Name:       "some-name",
								Type:       "some-custom-type",
								Source:     atc.Source{"some": "repository"},
								CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
							},
						},
					},
					pipeline.ConfigVersion(),
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeFalse())
			})

			It("returns the resource types tree given a resource", func() {
				resource, found, err := pipeline.Resource("some-name")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				tree := resourceTypes.Filter(resource)
				Expect(len(tree)).To(Equal(1))

				Expect(tree[0].Name()).To(Equal("some-name"))
				Expect(tree[0].Type()).To(Equal("some-custom-type"))
			})
		})

		Context("when building the dependency tree from resource types list", func() {
			BeforeEach(func() {
				var (
					created bool
					err     error
				)

				otherPipeline, created, err := defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-duplicate-type-name"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{
							{
								Name:       "some-custom-type",
								Type:       "some-different-foo-type",
								Source:     atc.Source{"some": "repository"},
								CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
							},
						},
					},
					db.ConfigVersion(0),
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeTrue())
				Expect(otherPipeline).NotTo(BeNil())

				pipeline, created, err = defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-types"},
					atc.Config{
						Resources: atc.ResourceConfigs{
							{
								Name:   "some-resource",
								Type:   "some-custom-type",
								Source: atc.Source{},
							},
						},
						ResourceTypes: atc.ResourceTypes{
							{
								Name:   "registry-image",
								Type:   "registry-image",
								Source: atc.Source{"some": "repository"},
							},
							{
								Name:       "some-other-type",
								Type:       "registry-image",
								Privileged: true,
								Source:     atc.Source{"some": "other-repository"},
							},
							{
								Name:   "some-type-with-params",
								Type:   "s3",
								Source: atc.Source{"some": "repository"},
								Params: atc.Params{"unpack": "true"},
							},
							{
								Name:       "some-type-with-custom-check",
								Type:       "registry-image",
								Source:     atc.Source{"some": "repository"},
								CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
							},
							{
								Name:       "some-custom-type",
								Type:       "some-other-foo-type",
								Source:     atc.Source{"some": "repository"},
								CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
							},
							{
								Name:       "some-other-foo-type",
								Type:       "some-other-type",
								Source:     atc.Source{"some": "repository"},
								CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
							},
						},
					},
					pipeline.ConfigVersion(),
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeFalse())
			})

			It("returns the resource types tree given a resource", func() {
				resource, found, err := pipeline.Resource("some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				tree := resourceTypes.Filter(resource)
				Expect(len(tree)).To(Equal(4))

				Expect(tree[0].Name()).To(Equal("some-custom-type"))
				Expect(tree[0].Type()).To(Equal("some-other-foo-type"))

				Expect(tree[1].Name()).To(Equal("some-other-foo-type"))
				Expect(tree[1].Type()).To(Equal("some-other-type"))

				Expect(tree[2].Name()).To(Equal("some-other-type"))
				Expect(tree[2].Type()).To(Equal("registry-image"))

				Expect(tree[3].Name()).To(Equal("registry-image"))
				Expect(tree[3].Type()).To(Equal("registry-image"))
			})
		})

		Context("Deserialize", func() {
			var rts atc.ResourceTypes
			JustBeforeEach(func() {
				rts = resourceTypes.Deserialize()
			})

			Context("when no base resource type defaults defined", func() {
				It("should return original resource types", func() {
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name:     "some-type",
						Type:     "registry-image",
						Source:   atc.Source{"some": "repository"},
						Defaults: atc.Source{"some-default-k1": "some-default-v1"},
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name: "some-other-type",
						Type: "some-type",
						Source: atc.Source{
							"some-default-k1": "some-default-v1",
							"some":            "other-repository",
						},
						Privileged: true,
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name:       "some-type-with-params",
						Type:       "s3",
						Source:     atc.Source{"some": "repository"},
						Defaults:   nil,
						Privileged: false,
						Params:     atc.Params{"unpack": "true"},
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name:       "some-type-with-custom-check",
						Type:       "registry-image",
						Source:     atc.Source{"some": "repository"},
						CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
					}))
				})
			})

			Context("when base resource type defaults is defined", func() {
				BeforeEach(func() {
					atc.LoadBaseResourceTypeDefaults(map[string]atc.Source{"s3": {"default-s3-key": "some-value"}})
				})
				AfterEach(func() {
					atc.LoadBaseResourceTypeDefaults(map[string]atc.Source{})
				})

				It("should return original resource types", func() {
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name:     "some-type",
						Type:     "registry-image",
						Source:   atc.Source{"some": "repository"},
						Defaults: atc.Source{"some-default-k1": "some-default-v1"},
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name: "some-other-type",
						Type: "some-type",
						Source: atc.Source{
							"some-default-k1": "some-default-v1",
							"some":            "other-repository",
						},
						Privileged: true,
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name: "some-type-with-params",
						Type: "s3",
						Source: atc.Source{
							"some":           "repository",
							"default-s3-key": "some-value",
						},
						Defaults:   nil,
						Privileged: false,
						Params:     atc.Params{"unpack": "true"},
					}))
					Expect(rts).To(ContainElement(atc.ResourceType{
						Name:       "some-type-with-custom-check",
						Type:       "registry-image",
						Source:     atc.Source{"some": "repository"},
						CheckEvery: &atc.CheckEvery{Interval: 10 * time.Millisecond},
					}))
				})
			})
		})
	})

	Describe("SetResourceConfigScope", func() {
		var resourceType db.ResourceType
		var scope db.ResourceConfigScope

		BeforeEach(func() {
			resourceType = defaultResourceType

			resourceConfig, err := resourceConfigFactory.FindOrCreateResourceConfig(resourceType.Type(), resourceType.Source(), nil)
			Expect(err).ToNot(HaveOccurred())

			scope, err = resourceConfig.FindOrCreateScope(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("associates the resource to the config and scope", func() {
			Expect(resourceType.ResourceConfigScopeID()).To(BeZero())

			Expect(resourceType.SetResourceConfigScope(scope)).To(Succeed())

			_, err := resourceType.Reload()
			Expect(err).ToNot(HaveOccurred())

			Expect(resourceType.ResourceConfigScopeID()).To(Equal(scope.ID()))
		})
	})

	Describe("CreateBuild", func() {
		var resourceType db.ResourceType
		var ctx context.Context
		var manuallyTriggered bool
		var plan atc.Plan

		var build db.Build
		var created bool

		BeforeEach(func() {
			ctx = context.TODO()
			resourceType = defaultResourceType
			manuallyTriggered = false
			plan = atc.Plan{
				ID: "some-plan",
				Check: &atc.CheckPlan{
					Name: "wreck",
				},
			}
		})

		JustBeforeEach(func() {
			var err error
			build, created, err = resourceType.CreateBuild(ctx, manuallyTriggered, plan)
			Expect(err).ToNot(HaveOccurred())
		})

		It("creates a started build for a resource type", func() {
			Expect(created).To(BeTrue())
			Expect(build).ToNot(BeNil())
			Expect(build.Name()).To(Equal(db.CheckBuildName))
			Expect(build.ResourceTypeID()).To(Equal(resourceType.ID()))
			Expect(build.PipelineID()).To(Equal(defaultResource.PipelineID()))
			Expect(build.TeamID()).To(Equal(defaultResource.TeamID()))
			Expect(build.IsManuallyTriggered()).To(BeFalse())
			Expect(build.Status()).To(Equal(db.BuildStatusStarted))
			Expect(build.PrivatePlan()).To(Equal(plan))
		})

		It("logs to the check_build_events partition", func() {
			err := build.SaveEvent(event.Log{Payload: "log"})
			Expect(err).ToNot(HaveOccurred())
			// created + log events
			Expect(numBuildEventsForCheck(build)).To(Equal(2))
		})

		Context("when tracing is configured", func() {
			var span trace.Span

			BeforeEach(func() {
				tracing.ConfigureTraceProvider(oteltest.NewTracerProvider())

				ctx, span = tracing.StartSpan(context.Background(), "fake-operation", nil)
			})

			AfterEach(func() {
				tracing.Configured = false
			})

			It("propagates span context", func() {
				traceID := span.SpanContext().TraceID().String()
				buildContext := build.SpanContext()
				traceParent := buildContext.Get("traceparent")
				Expect(traceParent).To(ContainSubstring(traceID))
			})
		})

		Context("when another build already exists", func() {
			var prevBuild db.Build

			BeforeEach(func() {
				var err error
				var prevCreated bool
				prevBuild, prevCreated, err = resourceType.CreateBuild(ctx, false, plan)
				Expect(err).ToNot(HaveOccurred())
				Expect(prevCreated).To(BeTrue())
			})

			It("does not create the second build", func() {
				Expect(created).To(BeFalse())
			})

			Context("when manually triggered", func() {
				BeforeEach(func() {
					manuallyTriggered = true
				})

				It("creates a manually triggered resource build", func() {
					Expect(created).To(BeTrue())
					Expect(build.IsManuallyTriggered()).To(BeTrue())
					Expect(build.ResourceTypeID()).To(Equal(resourceType.ID()))
				})
			})

			Context("when the previous build is finished", func() {
				BeforeEach(func() {
					Expect(prevBuild.Finish(db.BuildStatusSucceeded)).To(Succeed())
				})

				It("creates the build", func() {
					Expect(created).To(BeTrue())
					Expect(build.ResourceTypeID()).To(Equal(resourceType.ID()))
				})
			})
		})
	})

	Describe("CheckPlan", func() {

		Context("when there is a resource type using a base type", func() {
			var createdCheckPlan atc.Plan
			var version atc.Version
			var resourceType db.ResourceType

			BeforeEach(func() {
				pipeline, created, err := defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-resource-base-type"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{{
							Name: "some-resource-type",
							Type: "some-base-resource-type",
							Tags: []string{"tag"},
							Source: atc.Source{
								"some": "source",
							},
						}},
					},
					0,
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeTrue())

				var found bool
				resourceType, found, err = pipeline.ResourceType("some-resource-type")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())

				planFactory := atc.NewPlanFactory(0)
				version = atc.Version{"version": "from"}
				sourceDefault := atc.Source{"source-test": "default"}
				createdCheckPlan = resourceType.CheckPlan(planFactory, atc.ResourceTypes{}, version, 1*time.Hour, sourceDefault, false, false)
			})

			It("produces a simple check plan", func() {
				expectedPlan := atc.Plan{
					ID: atc.PlanID("1"),
					Check: &atc.CheckPlan{
						Name: resourceType.Name(),
						Type: resourceType.Type(),
						Source: atc.Source{
							"some":        "source",
							"source-test": "default",
						},
						Tags: resourceType.Tags(),
						TypeImage: atc.TypeImage{
							BaseType: resourceType.Type(),
						},
						FromVersion:  version,
						ResourceType: resourceType.Name(),
						Interval:     "1h0m0s",
					},
				}
				Expect(createdCheckPlan).To(Equal(expectedPlan))
			})
		})

		Context("when there is a resource type using a custom type", func() {

			var createdCheckPlan atc.Plan
			var version atc.Version
			var resourceType db.ResourceType

			BeforeEach(func() {
				pipeline, created, err := defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-resource-custom-type"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{
							{
								Name: "some-custom-resource-type",
								Type: "some-resource-type",
								Tags: []string{"tag"},
								Source: atc.Source{
									"some": "source",
								},
							},
							{
								Name:   "some-resource-type",
								Type:   "some-base-resource-type",
								Source: atc.Source{"some": "type-source"},
							},
						},
					},
					0,
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeTrue())

				var found bool
				resourceType, found, err = pipeline.ResourceType("some-custom-resource-type")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())

				planFactory := atc.NewPlanFactory(0)
				version = atc.Version{"version": "from"}
				createdCheckPlan = resourceType.CheckPlan(planFactory,
					atc.ResourceTypes{
						{
							Name:   "some-resource-type",
							Type:   "some-base-resource-type",
							Source: atc.Source{"some": "type-source"},
						},
					}, version, 1*time.Hour, nil, false, false)
			})

			It("produces a check plan with nested image steps", func() {
				checkPlanID := atc.PlanID("1/image-check")
				expectedPlan := atc.Plan{
					ID: atc.PlanID("1"),
					Check: &atc.CheckPlan{
						Name: resourceType.Name(),
						Type: resourceType.Type(),
						Source: atc.Source{
							"some": "source",
						},
						Tags: resourceType.Tags(),
						TypeImage: atc.TypeImage{
							BaseType: "some-base-resource-type",
							CheckPlan: &atc.Plan{
								ID: checkPlanID,
								Check: &atc.CheckPlan{
									Name:   "some-resource-type",
									Type:   "some-base-resource-type",
									Source: atc.Source{"some": "type-source"},
									TypeImage: atc.TypeImage{
										BaseType: "some-base-resource-type",
									},
									Tags: resourceType.Tags(),
								},
							},
							GetPlan: &atc.Plan{
								ID: atc.PlanID("1/image-get"),
								Get: &atc.GetPlan{
									Name:   "some-resource-type",
									Type:   "some-base-resource-type",
									Source: atc.Source{"some": "type-source"},
									TypeImage: atc.TypeImage{
										BaseType: "some-base-resource-type",
									},
									Tags:        resourceType.Tags(),
									VersionFrom: &checkPlanID,
								},
							},
						},
						FromVersion:  version,
						ResourceType: resourceType.Name(),
						Interval:     "1h0m0s",
					},
				}
				Expect(createdCheckPlan).To(Equal(expectedPlan))
			})
		})

		Context("when there is a resource type using a privileged custom type", func() {
			var createdCheckPlan atc.Plan
			var version atc.Version
			var resourceType db.ResourceType

			BeforeEach(func() {
				pipeline, created, err := defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-resource-custom-type"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{
							{
								Name: "some-custom-resource-type",
								Type: "some-resource-type",
								Source: atc.Source{
									"some": "source",
								},
							},
							{
								Name:       "some-resource-type",
								Type:       "some-base-resource-type",
								Source:     atc.Source{"some": "type-source"},
								Privileged: true,
							},
						},
					},
					0,
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeTrue())

				var found bool
				resourceType, found, err = pipeline.ResourceType("some-custom-resource-type")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())

				planFactory := atc.NewPlanFactory(0)
				version = atc.Version{"version": "from"}
				createdCheckPlan = resourceType.CheckPlan(planFactory,
					atc.ResourceTypes{
						{
							Name:       "some-resource-type",
							Type:       "some-base-resource-type",
							Source:     atc.Source{"some": "type-source"},
							Privileged: true,
						},
					}, version, 1*time.Hour, nil, false, false)
			})

			It("produces a check plan with privileged", func() {
				checkPlanID := atc.PlanID("1/image-check")
				expectedPlan := atc.Plan{
					ID: atc.PlanID("1"),
					Check: &atc.CheckPlan{
						Name: resourceType.Name(),
						Type: resourceType.Type(),
						Source: atc.Source{
							"some": "source",
						},
						TypeImage: atc.TypeImage{
							BaseType:   "some-base-resource-type",
							Privileged: true,
							CheckPlan: &atc.Plan{
								ID: checkPlanID,
								Check: &atc.CheckPlan{
									Name:   "some-resource-type",
									Type:   "some-base-resource-type",
									Source: atc.Source{"some": "type-source"},
									TypeImage: atc.TypeImage{
										BaseType: "some-base-resource-type",
									},
								},
							},
							GetPlan: &atc.Plan{
								ID: atc.PlanID("1/image-get"),
								Get: &atc.GetPlan{
									Name:   "some-resource-type",
									Type:   "some-base-resource-type",
									Source: atc.Source{"some": "type-source"},
									TypeImage: atc.TypeImage{
										BaseType: "some-base-resource-type",
									},
									VersionFrom: &checkPlanID,
								},
							},
						},
						FromVersion:  version,
						ResourceType: resourceType.Name(),
						Interval:     "1h0m0s",
					},
				}
				Expect(createdCheckPlan).To(Equal(expectedPlan))
			})
		})

		Context("when skipping the interval", func() {
			var resourceType db.ResourceType

			BeforeEach(func() {
				pipeline, created, err := defaultTeam.SavePipeline(
					atc.PipelineRef{Name: "pipeline-with-resource-custom-type"},
					atc.Config{
						ResourceTypes: atc.ResourceTypes{
							{
								Name: "some-custom-resource-type",
								Type: "some-resource-type",
								Source: atc.Source{
									"some": "source",
								},
							},
							{
								Name:   "some-resource-type",
								Type:   "some-base-resource-type",
								Source: atc.Source{"some": "type-source"},
							},
						},
					},
					0,
					false,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(BeTrue())

				var found bool
				resourceType, found, err = pipeline.ResourceType("some-custom-resource-type")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())
			})

			Context("when not skipping the interval recursively", func() {
				It("skips the interval for the resource check, but not for the resource type", func() {
					planFactory := atc.NewPlanFactory(0)
					version := atc.Version{"version": "from"}
					createdCheckPlan := resourceType.CheckPlan(planFactory,
						atc.ResourceTypes{
							{
								Name:   "some-resource-type",
								Type:   "some-base-resource-type",
								Source: atc.Source{"some": "type-source"},
							},
						}, version, 1*time.Hour, nil, true, false)

					checkPlanID := atc.PlanID("1/image-check")
					expectedPlan := atc.Plan{
						ID: atc.PlanID("1"),
						Check: &atc.CheckPlan{
							Name: resourceType.Name(),
							Type: resourceType.Type(),
							Source: atc.Source{
								"some": "source",
							},
							SkipInterval: true,
							TypeImage: atc.TypeImage{
								BaseType: "some-base-resource-type",
								CheckPlan: &atc.Plan{
									ID: checkPlanID,
									Check: &atc.CheckPlan{
										Name:         "some-resource-type",
										Type:         "some-base-resource-type",
										Source:       atc.Source{"some": "type-source"},
										SkipInterval: false,
										TypeImage: atc.TypeImage{
											BaseType: "some-base-resource-type",
										},
									},
								},
								GetPlan: &atc.Plan{
									ID: atc.PlanID("1/image-get"),
									Get: &atc.GetPlan{
										Name:   "some-resource-type",
										Type:   "some-base-resource-type",
										Source: atc.Source{"some": "type-source"},
										TypeImage: atc.TypeImage{
											BaseType: "some-base-resource-type",
										},
										VersionFrom: &checkPlanID,
									},
								},
							},
							FromVersion:  version,
							ResourceType: resourceType.Name(),
							Interval:     "1h0m0s",
						},
					}
					Expect(createdCheckPlan).To(Equal(expectedPlan))
				})
			})

			Context("when skipping the interval recursively", func() {
				It("skips the interval for the resource and resource type checks", func() {
					planFactory := atc.NewPlanFactory(0)
					version := atc.Version{"version": "from"}
					createdCheckPlan := resourceType.CheckPlan(planFactory,
						atc.ResourceTypes{
							{
								Name:   "some-resource-type",
								Type:   "some-base-resource-type",
								Source: atc.Source{"some": "type-source"},
							},
						}, version, 1*time.Hour, nil, true, true)

					checkPlanID := atc.PlanID("1/image-check")
					expectedPlan := atc.Plan{
						ID: atc.PlanID("1"),
						Check: &atc.CheckPlan{
							Name: resourceType.Name(),
							Type: resourceType.Type(),
							Source: atc.Source{
								"some": "source",
							},
							SkipInterval: true,
							TypeImage: atc.TypeImage{
								BaseType: "some-base-resource-type",
								CheckPlan: &atc.Plan{
									ID: checkPlanID,
									Check: &atc.CheckPlan{
										Name:         "some-resource-type",
										Type:         "some-base-resource-type",
										Source:       atc.Source{"some": "type-source"},
										SkipInterval: true,
										TypeImage: atc.TypeImage{
											BaseType: "some-base-resource-type",
										},
									},
								},
								GetPlan: &atc.Plan{
									ID: atc.PlanID("1/image-get"),
									Get: &atc.GetPlan{
										Name:   "some-resource-type",
										Type:   "some-base-resource-type",
										Source: atc.Source{"some": "type-source"},
										TypeImage: atc.TypeImage{
											BaseType: "some-base-resource-type",
										},
										VersionFrom: &checkPlanID,
									},
								},
							},
							FromVersion:  version,
							ResourceType: resourceType.Name(),
							Interval:     "1h0m0s",
						},
					}
					Expect(createdCheckPlan).To(Equal(expectedPlan))
				})
			})
		})
	})
})
